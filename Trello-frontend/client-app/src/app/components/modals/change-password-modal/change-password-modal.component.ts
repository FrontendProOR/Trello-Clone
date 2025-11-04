import { Component } from '@angular/core';
import { ReactiveFormsModule, FormGroup, FormBuilder, Validators } from '@angular/forms';
import { CommonModule } from '@angular/common';  // Import CommonModule
import { NgbActiveModal } from '@ng-bootstrap/ng-bootstrap';
import { AuthService } from '../../../services/auth.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-change-password-modal',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],  // Add CommonModule here
  templateUrl: './change-password-modal.component.html',
  styleUrls: ['./change-password-modal.component.scss']
})
export class ChangePasswordModalComponent {
  passwordForm: FormGroup;
  passwordMismatch: boolean = false;

  constructor(public activeModal: NgbActiveModal, private fb: FormBuilder, private authService: AuthService, private router: Router) {
    this.passwordForm = this.fb.group({
      oldPassword: ['', Validators.required],
      newPassword: ['', [Validators.required, Validators.minLength(6)]],
      confirmPassword: ['', Validators.required]
    });
  }

  onSubmit() {
    const token = this.authService.getToken();
    let email = '';
    
    if (token) {
      const secondPartOfToken = (token || "").split('.')[1];
      const decodedToken = JSON.parse(atob(secondPartOfToken));
      email = decodedToken.email;
    }

    if (this.passwordForm.invalid) {
      return;
    }

    const { oldPassword, newPassword, confirmPassword } = this.passwordForm.value;

    if (newPassword !== confirmPassword) {
      this.passwordMismatch = true;
      return;
    }

    this.authService.changePassword(email, oldPassword, newPassword).subscribe(
      (response) => {
        this.activeModal.close('Password changed successfully');
        alert('Password changed successfully');
        this.authService.logout();
        this.activeModal.close('Password changed successfully');
        this.router.navigate(['/users/login']);

      },
      (error) => {
        console.error('Error changing password:', error);
        this.passwordMismatch = false;
        alert('Error changing password: ' + error.error?.message || 'An unknown error occurred');
      }
    );
  }


  closeModal() {
    this.activeModal.dismiss('Cancelled');
  }
}
