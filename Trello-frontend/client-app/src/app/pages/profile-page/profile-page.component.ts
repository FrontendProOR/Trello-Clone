import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { AuthService } from '../../services/auth.service';
import { UserService } from '../../services/user.service';
import { User } from '../../models/user.model';
import { ProjectService } from '../../services/project.service';
import { CommonModule } from '@angular/common';
import { ChangePasswordModalComponent } from '../../components/modals/change-password-modal/change-password-modal.component';
import { NgbModal } from '@ng-bootstrap/ng-bootstrap';
import { DeleteAccountModalComponent } from '../../components/modals/delete-account-modal/delete-account-modal.component';
import { EditUserModalComponent } from '../../components/modals/edit-user-modal/edit-user-modal.component';

@Component({
  selector: 'app-profile-page',
  standalone: true,
  imports: [CommonModule, ChangePasswordModalComponent],
  templateUrl: './profile-page.component.html',
  styleUrls: ['./profile-page.component.scss'],
})
export class ProfilePageComponent implements OnInit {
  user: User | null = null;
  role: string = '';
  managerProjects: any[] = [];
  memberProjects: any[] = [];
  showChangePasswordModal: boolean = false;

  constructor(
    private authService: AuthService,
    private userService: UserService,
    private projectService: ProjectService,
    private modalService: NgbModal,
    private router: Router
  ) {}

  ngOnInit(): void {
    if (!this.authService.isLoggedIn()) {
      this.router.navigate(['users/login']);
      return;
    }

    const token = this.authService.getToken();
    if (!token) {
      this.router.navigate(['users/login']);
      return;
    }

    const secondPartOfToken = (token || '').split('.')[1];
    if (!secondPartOfToken) {
      this.router.navigate(['users/login']);
      return;
    }

    const decodedPayload = atob(secondPartOfToken);
    if (!decodedPayload) {
      this.router.navigate(['users/login']);
      return;
    }

    const decodedJson = JSON.parse(decodedPayload);
    if (!decodedJson || !decodedJson.user_id) {
      this.router.navigate(['users/login']);
      return;
    }

    const userId = decodedJson.user_id;
    console.log('User ID:', userId);
    this.loadUserProfile(userId);
  }

  loadUserProfile(userId: string): void {
    this.userService.getUserById(userId).subscribe(
      (user) => {
        this.user = user;
        // Direktno unescaping polja pre prikaza
        this.user.first_name = this.unescapeHTML(user.first_name);
        this.user.last_name = this.unescapeHTML(user.last_name);
        this.user.username = this.unescapeHTML(user.username);
        this.user.email = this.unescapeHTML(user.email);
        this.role = user.role;
        if (this.role === 'manager') {
          this.loadManagerProjects();
        } else {
          this.loadMemberProjects();
        }
      },
      (error) => {
        console.error('Error loading user profile:', error);
      }
    );
  }

  unescapeHTML(html: string): string {
    return html
      .replace(/&lt;/g, '<')
      .replace(/&gt;/g, '>')
      .replace(/&amp;/g, '&')
      .replace(/&quot;/g, '"')
      .replace(/&#39;/g, "'");
  }

  loadManagerProjects(): void {
    this.projectService.getProjectsByManagerId(this.user!.id).subscribe(
      (projects) => {
        this.managerProjects = projects;
      },
      (error) => {
        console.error('Error loading manager projects:', error);
      }
    );
  }

  loadMemberProjects(): void {
    this.projectService.getProjectsByMemberId(this.user!.id).subscribe(
      (projects) => {
        this.memberProjects = projects;
      },
      (error) => {
        console.error('Error loading member projects:', error);
      }
    );
  }

  openChangePasswordModal(projectId: string) {
    const modalRef = this.modalService.open(ChangePasswordModalComponent);
    modalRef.componentInstance.projectId = projectId;
  }

  closeChangePasswordModal() {
    this.modalService.dismissAll();
  }

  openDeleteAccountModal() {
    const modalRef = this.modalService.open(DeleteAccountModalComponent);

    modalRef.result.then(
      (result) => {
        if (result === 'confirmed') {
          this.deleteAccount();
        }
      },
      (reason) => {
        console.log('Modal dismissed', reason);
      }
    );
  }

  openEditUserModal(userId: string) {
    const modalRef = this.modalService.open(EditUserModalComponent);
    modalRef.result.then(
      (result) => {
        if (result === 'confirmed') {
          this.updateUser();
        }
      },
      (reason) => {
        console.log('Modal dismissed', reason);
      }
    );
  }

  deleteAccount() {
    if (this.user) {
      this.userService.deleteUser(this.user.id).subscribe(
        () => {
          this.authService.logout();
          this.router.navigate(['users/login']);
        },
        (error) => {
          console.error('Error deleting account:', error);
        }
      );
    }
  }

  updateUser() {
    this.router.navigate(['/users/edit', this.user?.id]);
  }
}
