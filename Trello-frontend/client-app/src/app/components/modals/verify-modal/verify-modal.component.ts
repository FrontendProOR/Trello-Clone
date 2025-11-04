import { Component } from '@angular/core';
import { NgbModal } from '@ng-bootstrap/ng-bootstrap';
import { FormsModule } from '@angular/forms';
import { AuthService } from '../../../services/auth.service';
import { LoginModalComponent } from '../login-modal/login-modal.component';

@Component({
  selector: 'app-verify-modal',
  standalone: true,
  imports: [FormsModule],
  templateUrl: './verify-modal.component.html',
  styleUrls: ['./verify-modal.component.scss']
})
export class VerifyModalComponent {
  email: string = '';
  code: string = '';

  constructor(private modalService: NgbModal, private authService: AuthService) {}

  onSubmit(): void {
    this.authService.verify(this.email, this.code).subscribe(
      (response) => {
        console.log('Verification successful:', response);
        this.sendMagicLink(); 
      },
      (error) => {
        console.error('Verification failed:', error);
        this.closeModal();
      }
    );
  }

  sendMagicLink(): void {
    this.authService.sendMagicLink(this.email).subscribe(
      (response) => {
        console.log('Magic link sent successfully:', response);
        this.openLoginModal();
        this.closeModal();
      },
      (error) => {
        console.error('Failed to send magic link:', error);
        this.closeModal();
      }
    );
  }

  submitVerification(): void {
    this.onSubmit();
  }

  openLoginModal() {
    this.modalService.open(LoginModalComponent);
  }

  closeModal() {
    this.modalService.dismissAll();
  }
}
