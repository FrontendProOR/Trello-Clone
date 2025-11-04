import { Component } from '@angular/core';
import { NgbModal } from '@ng-bootstrap/ng-bootstrap';
import { FormsModule } from '@angular/forms';
import { LoginModalComponent } from '../modals/login-modal/login-modal.component';
import { RegisterModalComponent } from '../modals/register-modal/register-modal.component';
import { VerifyModalComponent } from '../modals/verify-modal/verify-modal.component';
import { DashboardComponent } from '../../pages/dashboard/dashboard.component';
import { NotificationComponent } from './../notification/notification.component';
import { ActivityComponent } from './../activity/activity.component';
import { CommonModule } from '@angular/common';
import { AuthService } from '../../services/auth.service';

@Component({
  selector: 'app-navbar',
  standalone: true,
  imports: [FormsModule, RegisterModalComponent, LoginModalComponent, VerifyModalComponent,CommonModule],
  templateUrl: './navbar.component.html',
  styleUrls: ['./navbar.component.scss']
})

export class NavbarComponent {
  constructor(
    private modalService: NgbModal,
    public authService: AuthService 
  ) {}

  openRegisterModal() {
    this.modalService.open(RegisterModalComponent);
  }

  openLoginModal() {
    this.modalService.open(LoginModalComponent);
  }

  logout() {
    this.authService.logout();
    this.modalService.dismissAll();
    window.location.replace('/');
  }

  profile() {
    window.location.replace('users/profile');
  }

  dashboard() {
    window.location.replace('/dashboard');
  }

  notifications() {
    this.modalService.open(NotificationComponent);
  }

  activity() {
    this.modalService.open(ActivityComponent);
  }
  
}
