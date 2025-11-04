import { Component, Input } from '@angular/core';
import { UserService } from '../../../services/user.service';  // Import UserService
import { AuthService } from '../../../services/auth.service';  // Import AuthService
import { Router } from '@angular/router';
import { NgbActiveModal } from '@ng-bootstrap/ng-bootstrap';

@Component({
  selector: 'app-add-member-modal',
  templateUrl: './add-member-modal.component.html',
  styleUrls: ['./add-member-modal.component.scss']
})
export class AddMemberModalComponent {
  @Input() projectId!: string;  // Accept projectId dynamically

  constructor(private userService: UserService, private authService: AuthService, private router: Router, private activeModal: NgbActiveModal) {}

  // Method to add member to project by email
  addMember(email: string) {
    const token = this.authService.getToken();
    if (!token) {
      alert('You are not authorized to perform this action. Please log in.');
      return;
    }

    if (!email) {
      alert('Email is required.');
      return;
    }

    this.userService.addUserToProjectByEmail(this.projectId, email).subscribe({
      next: (response : any) => {
        alert('Member added successfully');
        const modal = document.getElementById('addMemberModal');
        if (modal) {
          modal?.classList.remove('show');
        }
      },
      error: (error: any) => {
        console.error('Error adding member', error);
        alert('Failed to add member. Please ensure the email is correct and try again.');
      }
    });
  }

  public closeModal() {
    this.activeModal.close();
  }
}
