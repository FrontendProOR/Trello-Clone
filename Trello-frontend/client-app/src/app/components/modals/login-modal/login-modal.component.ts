import { Component } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { NgbModal } from '@ng-bootstrap/ng-bootstrap';
import { AuthService } from '../../../services/auth.service';
import { Router } from '@angular/router';
import { HttpClientModule } from '@angular/common/http';
import { AppForgotPasswordModalComponent } from '../app-forgot-password-modal/app-forgot-password-modal.component';

@Component({
  selector: 'app-login-modal',
  standalone: true,
  imports: [FormsModule, HttpClientModule],
  templateUrl: './login-modal.component.html',
  styleUrls: ['./login-modal.component.scss']
})

export class LoginModalComponent {

  loginData = { email: '', password: '' };  
  errorMessage: string = '';

  constructor(private modalService: NgbModal, private authService: AuthService, private router: Router) {}

  onLoginSubmit(): void {  
    // console.log('Login data:', this.loginData);

    this.authService.login(this.loginData.email, this.loginData.password).subscribe(
      response => {
        localStorage.setItem('token', response.token);
        this.closeLoginModal();
        this.router.navigate(['/dashboard']);
      },
      error => {
        this.errorMessage = 'Invalid credentials. Please try again.';
      }
    );
  }

  openLoginModal() {
    this.modalService.open(this, {
      ariaLabelledBy: 'loginModalLabel',
      backdrop: 'static',  
      keyboard: false,     
    });
  }
  
  closeLoginModal() {
    this.modalService.dismissAll(); 
  }

  openRegisterModal() {
    this.closeLoginModal();
    this.modalService.open('app-register-modal', {
      ariaLabelledBy: 'registerModalLabel',
      backdrop: 'static',
      keyboard: false,
    });
  }

    openForgotPasswordModal() {
    this.closeLoginModal();
    this.modalService.open(AppForgotPasswordModalComponent, {
      ariaLabelledBy: 'forgotPasswordModalLabel',
      backdrop: 'static',
      keyboard: false,
    });
  }


}
