import { Component } from '@angular/core';
import { ProjectService } from '../../../services/project.service';
import { NgbModal } from '@ng-bootstrap/ng-bootstrap';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';
import { AuthService } from '../../../services/auth.service';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Router } from '@angular/router';

@Component({
  selector: 'app-add-project-modal',
  standalone: true,
  imports: [FormsModule, ReactiveFormsModule, CommonModule],
  templateUrl: './add-project-modal.component.html',
  styleUrls: ['./add-project-modal.component.scss']
})
export class AddProjectModalComponent {
  projectForm: FormGroup;

  constructor(
    private projectService: ProjectService,
    private modalService: NgbModal,
    private authService: AuthService,
    private fb: FormBuilder,
    private router: Router
  ) {
    this.projectForm = this.fb.group({
      name: ['', [Validators.required, Validators.maxLength(100)]], // max 100 characters
      description: ['', [Validators.required, Validators.maxLength(200)]], // max 200 characters
      expectedEndDate: ['', [Validators.required, this.expectedEndDateValidator]],
      minMembers: [2, [Validators.required, Validators.min(1), Validators.max(99)]],
      maxMembers: [5, [Validators.required, Validators.min(1), Validators.max(99)]],
      managerId: [''], 
      status: ['planning', [Validators.required]],
      workflowGraphId: ['new-workflow-graph-id'], 
      documents: [[]], 
      notificationsEnabled: [true], 
      history: [[]], 
      analytics: this.fb.group({
        taskCompletionRate: [0],
        overdueTasks: [0],
        activeMembers: [0]
      })
    }, { validator: this.minMaxValidator });
  }

  // Validator for checking minMembers and maxMembers
  minMaxValidator(group: FormGroup): {[key: string]: boolean} | null {
    const minMembers = group.get('minMembers')?.value;
    const maxMembers = group.get('maxMembers')?.value;
    return minMembers <= maxMembers ? null : { minMaxError: true };
  }

  // Custom validator for expectedEndDate
  expectedEndDateValidator(control: any): {[key: string]: boolean} | null {
    const today = new Date();
    const selectedDate = new Date(control.value);
    if (selectedDate < today) {
      return { dateBeforeToday: true };
    }
    const year2030 = new Date('2030-01-01');
    if (selectedDate > year2030) {
      return { dateAfter2030: true };
    }
    return null;
  }

  addProject() {
    const token = this.authService.getToken();
    if (!token) {
      this.closeAddProjectModal();
      this.router.navigate(['users/login']);
      return;
    }

    const secondPartOfToken = (token || "").split('.')[1];
    const decodedPayload = atob(secondPartOfToken);
    const decodedJson = JSON.parse(decodedPayload);
    if (!decodedJson || !decodedJson.user_id) {
      this.closeAddProjectModal();
      this.router.navigate(['users/login']);
      return; 
    }

    const userId = decodedJson.user_id;
    this.projectForm.patchValue({ managerId: userId });

    console.log('Form Data:', this.projectForm.value);
    console.log('User ID:', userId);

    // Ensure expectedEndDate is in correct format (yyyy-MM-dd'T'HH:mm:ss'Z')
    const formattedEndDate = this.formatDateForBackend(this.projectForm.value.expectedEndDate);
    this.projectForm.patchValue({ expectedEndDate: formattedEndDate });

    if (this.projectForm.valid) {
      const project = this.projectForm.value;
      this.projectService.addProject(project).subscribe(() => {
        window.location.reload();
        this.closeAddProjectModal();
      });
    }
  }

  openAddProjectModal() {
    this.modalService.open(this, {
      ariaLabelledBy: 'addProjectModalLabel',
      backdrop: 'static',
      keyboard: false
    });
  }

  closeAddProjectModal() {
    this.modalService.dismissAll();
  }

  // Helper method to format the date as "yyyy-MM-dd'T'HH:mm:ss'Z'"
  formatDateForBackend(date: string): string {
    const newDate = new Date(date);  // Convert to Date object
    return newDate.toISOString();    // Converts to ISO 8601 format (e.g., 2024-12-02T13:17:23.000Z)
  }
}
