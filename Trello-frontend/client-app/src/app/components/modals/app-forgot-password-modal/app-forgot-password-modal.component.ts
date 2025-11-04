import { Component } from '@angular/core';
import { NgbActiveModal } from '@ng-bootstrap/ng-bootstrap';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-app-forgot-password-modal',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './app-forgot-password-modal.component.html',
  styleUrls: ['./app-forgot-password-modal.component.scss']
})
export class AppForgotPasswordModalComponent {
  email: string = '';

  constructor(public activeModal: NgbActiveModal) {}

  submit() {
    if (!this.email) return;
    // Pozovi servis za slanje linka ili reset lozinke
    console.log('Reset password for:', this.email);
    this.activeModal.close();
  }
}


