import { Component, OnInit } from '@angular/core';
import { NgbActiveModal } from '@ng-bootstrap/ng-bootstrap';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';
import { AuthService } from '../../../services/auth.service';
import { UserService } from '../../../services/user.service';
import { User } from '../../../models/user.model';
import { CommonModule } from '@angular/common'; // For common directives (e.g., ngIf)

@Component({
  selector: 'app-edit-user-modal',
  standalone: true,
  imports: [ReactiveFormsModule, CommonModule], // Import ReactiveFormsModule here
  templateUrl: './edit-user-modal.component.html',
  styleUrls: ['./edit-user-modal.component.scss']
})
export class EditUserModalComponent implements OnInit {
  editUserForm: FormGroup;
  user: User | null = null;

  constructor(
    private fb: FormBuilder,
    public activeModal: NgbActiveModal,
    private authService: AuthService,
    private userService: UserService
  ) {
    this.editUserForm = this.fb.group({
      first_name: ['', Validators.required],
      last_name: ['', Validators.required],
      username: ['', [Validators.required, Validators.minLength(3)]],
      email: ['', [Validators.required, Validators.email]],
      role: ['', Validators.required]
    });
  }

  ngOnInit(): void {
    this.loadUserData();
  }

  // Load user data from the token
  loadUserData(): void {
    const token = this.authService.getToken();
    if (!token) {
      return;
    }

    const decodedPayload = this.decodeToken(token);
    if (decodedPayload?.user_id) {
      this.userService.getUserById(decodedPayload.user_id).subscribe(
        (user) => {
          this.user = user;
          this.editUserForm.patchValue(user);
        },
        (error) => {
          console.error('Error loading user data:', error);
        }
      );
    }
  }

  // Decode JWT token to get the payload
  decodeToken(token: string): any {
    const secondPartOfToken = token.split('.')[1];
    const decodedPayload = atob(secondPartOfToken);
    return JSON.parse(decodedPayload);
  }

  // Handle form submission to update user data
  onSubmit(): void {
    if (this.editUserForm.valid) {
      console.log('Updated Form Data:', this.editUserForm.value);
      console.log('User ID:', this.user!.id);
      this.userService.updateUser(this.user!.id, this.editUserForm.value).subscribe(
        (response) => {
          console.log('Updated Form Data:', this.editUserForm.value);
          console.log('User ID:', this.user!.id);
          console.log('User updated successfully:', response);
          this.activeModal.close();
        },
        (error) => {
          console.error('Error updating user:', error);
        }
      );
    } else {
      console.log('Form is invalid');
    }
  }

  // Close the modal
  close() {
    this.activeModal.dismiss();
  }
}
