import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';
import { HttpClientModule } from '@angular/common/http'; 
import { ProjectService } from '../../services/project.service'; 
import { NgbModal, NgbModalModule } from '@ng-bootstrap/ng-bootstrap';
import { AddProjectModalComponent } from '../modals/add-project-modal/add-project-modal.component';
import { EditProjectModalComponent } from '../modals/edit-project-modal/edit-project-modal.component';
import { ProjectPageComponent } from '../../pages/project-page/project-page.component';
import { AuthService } from '../../services/auth.service';
import { AddMemberModalComponent } from '../modals/add-member-modal/add-member-modal.component';
import { RemoveMemberModalComponent
  
 } from '../modals/remove-member-modal/remove-member-modal.component';
@Component({
  selector: 'app-project-list',
  standalone: true,
  imports: [CommonModule, HttpClientModule, NgbModalModule, AddProjectModalComponent], 
  templateUrl: './project-list.component.html',
  styleUrls: ['./project-list.component.scss'],
})
export class ProjectListComponent {
  managerProjects: any[] = [];
  memberProjects: any[] = [];
  userRole: string = '';

  constructor(
    private projectService: ProjectService, 
    private router: Router, 
    private modalService: NgbModal,
    private authService: AuthService
  ) {
    this.loadUserRoleAndProjects();
  }

  loadUserRoleAndProjects() {
    const token = this.authService.getToken();
    if (!token) {
      this.router.navigate(['users/login']);
      return;
    }

    const secondPartOfToken = (token || "").split('.')[1];
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
    this.userRole = decodedJson.role; // Assuming the role is stored in the decoded token

    // Load the projects based on user role
    if (this.userRole === 'manager') {
      this.loadProjectsForManagerSection(userId);
    } else if (this.userRole === 'member') {
      this.loadProjectsForMemberSection(userId);
    }
  }

  loadProjectsForManagerSection(userId: string) {
    this.projectService.getProjectsByManagerId(userId).subscribe((data) => {
      console.log("Manager Projects:", data);
      this.managerProjects = data;
    });
  }

  loadProjectsForMemberSection(userId: string) {
    this.projectService.getProjectsByMemberId(userId).subscribe((data) => {
      console.log("Member Projects:", data);
      this.memberProjects = data;
    });
  }

  goToProjectPage(projectId: string): void {
    this.router.navigate([`/projects/${projectId}/tasks`]);
  }

  formatDate(date: Date): string {
    return new Date(date).toLocaleDateString();
  }

  deleteProject(projectId: string): void {
    this.projectService.deleteProject(projectId).subscribe(() => {
      this.loadUserRoleAndProjects();
      this.router.navigate(['/dashboard']);
    });
  }

  openAddProjectModal() {
    this.modalService.open(AddProjectModalComponent);
  }

  closeAddProjectModal() {
    this.modalService.dismissAll();
  }

  openEditProjectModal(projectId: string) {
    const modalRef = this.modalService.open(EditProjectModalComponent);
    modalRef.componentInstance.projectId = projectId;
  }
  
  closeEditProjectModal() {
    this.modalService.dismissAll();
  }

  openAddMemberModal(projectId: string) {
    const modalRef = this.modalService.open(AddMemberModalComponent);
    modalRef.componentInstance.project = projectId;
  }

  closeAddMemberModal() {
    this.modalService.dismissAll();
  }

  openRemoveMemberModal(projectId: string) {
    const modalRef = this.modalService.open(RemoveMemberModalComponent);
    modalRef.componentInstance.project = projectId;
  }

  closeRemoveMemberModal() {
    this.modalService.dismissAll();
  }
}
