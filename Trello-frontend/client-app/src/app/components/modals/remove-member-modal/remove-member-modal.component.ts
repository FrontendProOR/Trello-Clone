import { Component, Input, OnInit } from '@angular/core';
import { UserService } from '../../../services/user.service'; // Import UserService
import { AuthService } from '../../../services/auth.service'; // Import AuthService
import { NgbActiveModal } from '@ng-bootstrap/ng-bootstrap'; // NgbActiveModal for modal close

@Component({
  selector: 'app-remove-member-modal',
  templateUrl: './remove-member-modal.component.html',
  styleUrls: ['./remove-member-modal.component.scss']
})
export class RemoveMemberModalComponent implements OnInit {
  @Input() projectId!: string; // Accept projectId dynamically
  members: any[] = []; // Array to hold the project members

  constructor(
    private userService: UserService,
    private authService: AuthService,
    public activeModal: NgbActiveModal
  ) {}

  ngOnInit() {
    this.loadMembers(); // Load members on initialization
  }

  loadMembers() {
    // Fetch members from the backend
    this.userService.getAllMembersInProject(this.projectId).subscribe({
      next: (members : any) => {
        this.members = members; // Store the members in the component
      },
      error: (err : any) => {
        console.error('Error fetching members:', err);
        alert('Failed to load project members.');
      }
    });
  }

  removeMember(userId: string) {
    const token = this.authService.getToken();
    if (!token) {
      alert('You are not authorized to perform this action. Please log in.');
      return;
    }

    this.userService.removeUserFromProject(this.projectId, userId).subscribe({
      next: () => {
        alert('Member removed successfully.');
        this.loadMembers(); // Reload the list of members
      },
      error: (error: any) => {
        console.error('Error removing member', error);
        alert('Failed to remove member.');
      }
    });
  }
}
