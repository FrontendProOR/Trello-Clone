import { Component } from '@angular/core';
import { RouterModule } from '@angular/router';
import { LoginModalComponent } from '../../components/modals/login-modal/login-modal.component';
import { RegisterModalComponent } from '../../components/modals/register-modal/register-modal.component';
import { FormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';
import { VerifyModalComponent } from '../../components/modals/verify-modal/verify-modal.component';
import { NgbModal } from '@ng-bootstrap/ng-bootstrap';

@Component({
  selector: 'app-homepage',
  standalone: true,
  imports: [RouterModule,LoginModalComponent,FormsModule,CommonModule,RegisterModalComponent,VerifyModalComponent],
  templateUrl: './homepage.component.html',
  styleUrl: './homepage.component.scss',
  template: `
    <div class="container">
      <h1>Welcome to the Trello Clone</h1>
      <p>Please login or register to continue.</p>
    </div>
  `,
  styles: [
    `
      .container {
        text-align: center;
        margin-top: 50px;
      }
    `,
  ]
})
export class HomepageComponent {

  constructor(private modalService: NgbModal) {}

  openLoginModal() {
    this.modalService.open(LoginModalComponent);
  }

  openRegisterModal() {
    this.modalService.open(RegisterModalComponent);
  }
}
