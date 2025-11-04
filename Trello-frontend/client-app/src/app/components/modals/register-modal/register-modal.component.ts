import { Component } from '@angular/core';
import { NgbModal } from '@ng-bootstrap/ng-bootstrap';
import { Router } from '@angular/router';
import { debounceTime, distinctUntilChanged, switchMap } from 'rxjs/operators';
import { of } from 'rxjs';

import { FormsModule } from '@angular/forms';
import { AuthService } from '../../../services/auth.service';
import {
  FormBuilder,
  FormGroup,
  Validators,
  ReactiveFormsModule,
} from '@angular/forms';
import { CommonModule } from '@angular/common';
import { VerifyModalComponent } from '../verify-modal/verify-modal.component';

@Component({
  selector: 'app-register-modal',
  standalone: true,
  imports: [FormsModule, ReactiveFormsModule, CommonModule],
  templateUrl: './register-modal.component.html',
  styleUrls: ['./register-modal.component.scss'],
})
export class RegisterModalComponent {
  registerForm: FormGroup;
  passwordBlacklistError: boolean = false;
  captchaImageUrl: string = '';
  captchaId: string = '';

  constructor(
    private fb: FormBuilder,
    private modalService: NgbModal,
    private router: Router,
    private authService: AuthService
  ) {
    this.registerForm = this.fb.group(
      {
        first_name: ['', Validators.required],
        last_name: ['', Validators.required],
        username: ['', [Validators.required, Validators.minLength(3)]],
        email: ['', [Validators.required, Validators.email]],
        password: ['', [Validators.required, Validators.minLength(6)]],
        confirm_password: ['', [Validators.required]],
        role: ['', Validators.required],
        captchaText: ['', Validators.required],
      },
      {
        validator: this.passwordMatchValidator,
      }
    );

    // Handle password blacklist checks
    this.registerForm.controls['password'].valueChanges
      .pipe(
        debounceTime(500),
        distinctUntilChanged(),
        switchMap((password) =>
          password.length >= 6
            ? this.authService.checkPasswordBlacklist(password)
            : of(false)
        )
      )
      .subscribe((isBlacklisted) => {
        this.passwordBlacklistError = isBlacklisted;
        if (isBlacklisted) {
          this.registerForm.controls['password'].setErrors({
            blacklisted: true,
          });
        } else {
          this.registerForm.controls['password'].setErrors(null);
        }
      });

    // Generate initial CAPTCHA
    this.generateCaptcha();
  }

  passwordMatchValidator(group: FormGroup): { [key: string]: boolean } | null {
    const password = group.get('password')?.value;
    const confirmPassword = group.get('confirm_password')?.value;
    return password === confirmPassword ? null : { passwordsMismatch: true };
  }

  generateCaptcha(): void {
    this.authService.generateCaptcha().subscribe({
      next: (response: { id: string; image: string }) => {
        this.captchaId = response.id; // ID CAPTCHA
        this.captchaImageUrl = response.image; // Base64 slika
      },
      error: () => {
        alert('Failed to generate CAPTCHA. Please try again.');
      },
    });
  }

  // Funkcija za osvežavanje CAPTCHA
  refreshCaptcha(): void {
    this.generateCaptcha();
  }

  onSubmit(): void {
    if (this.registerForm.valid) {
      const captchaId = this.captchaId;
      const captchaText = this.registerForm.value.captchaText;

      console.log('Sending CAPTCHA ID:', captchaId);
      console.log('Sending CAPTCHA Text:', captchaText);

      this.authService.validateCaptcha(captchaId, captchaText).subscribe({
        next: () => {
          // CAPTCHA validna, nastavite sa registracijom
          const formData = this.registerForm.value;
          this.authService.registerWithCaptcha(formData).subscribe(
            (response) => {
              console.log('CAPTCHA validated successfully!', response);

              console.log('Registration successful:', response);
              this.closeRegisterModal();
              this.openVerifyModal();
            },
            (error) => {
              console.error('Registration failed:', error);
              this.refreshCaptcha(); // Osvežavanje CAPTCHA
            }
          );
        },
        error: (err) => {
          console.error('CAPTCHA validation failed:', err);
          alert('Invalid CAPTCHA. Please try again.');
          this.refreshCaptcha(); // Osvežavanje CAPTCHA
        },
      });
    } else {
      console.log('Form is invalid');
    }
  }

  openVerifyModal() {
    this.modalService.open(VerifyModalComponent);
  }

  openRegisterModal() {
    this.modalService.open(this, {
      ariaLabelledBy: 'registerModalLabel',
      backdrop: 'static',
      keyboard: false,
    });
  }

  closeRegisterModal() {
    this.modalService.dismissAll();
  }
}
