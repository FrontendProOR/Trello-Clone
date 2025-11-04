import { Component, OnInit, Input } from '@angular/core';
import { ProjectService } from '../../../services/project.service';
import { NgbModal } from '@ng-bootstrap/ng-bootstrap';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { ActivatedRoute } from '@angular/router';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule } from '@angular/forms';
import { FormsModule } from '@angular/forms';
import { HttpClientModule } from '@angular/common/http';

@Component({
  selector: 'app-edit-project-modal',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, FormsModule, HttpClientModule],
  templateUrl: './edit-project-modal.component.html',
  styleUrls: ['./edit-project-modal.component.scss']
})
export class EditProjectModalComponent implements OnInit {
  @Input() projectId: string = '';  // Use Input to receive projectId
  projectForm: FormGroup;

  constructor(
    private projectService: ProjectService,
    private modalService: NgbModal,
    private fb: FormBuilder,
  ) { 
    this.projectForm = this.fb.group({
      name: ['', [Validators.required, Validators.maxLength(100)]], // max 100 characters
      description: ['', [Validators.required, Validators.maxLength(200)]], // max 200 characters
      expectedEndDate: ['', [Validators.required, this.expectedEndDateValidator]],
      createdAt: [''],
      minMembers: [2, [Validators.required, Validators.min(1), Validators.max(99)]],
      maxMembers: [5, [Validators.required, Validators.min(1), Validators.max(99)]],
      status: ['planning', [Validators.required]],
      managerId: [''], 
      workflowGraphId: [''],
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

  ngOnInit(): void {
    this.loadProjectData();
  }

  loadProjectData(): void {
    if (this.projectId) {
      this.projectService.getProjectById(this.projectId).subscribe(
        (project) => {
          if (project) {
            const formattedEndDate = this.formatDateForInput(project.expectedEndDate);
            this.projectForm.patchValue({
              name: project.name,
              description: project.description,
              expectedEndDate: formattedEndDate,
              createdAt: project.createdAt,
              minMembers: project.minMembers,
              maxMembers: project.maxMembers,
              status: project.status,
              managerId: project.managerId,
              workflowGraphId: project.workflowGraphId,
              documents: project.documents,
              notificationsEnabled: project.notificationsEnabled,
              history: project.history,
              analytics: project.analytics
            });
          }
        },
        error => {
          console.error('Error loading project details', error);
        }
      );
    }
  }

  formatDateForInput(date: string): string {
    const newDate = new Date(date);
    return newDate.toISOString().split('T')[0];
  }

  editProject(): void {
    if (this.projectForm.invalid) {
      return;
    }

    const projectData = this.projectForm.value;
    if (projectData.expectedEndDate) {
      projectData.expectedEndDate = new Date(projectData.expectedEndDate).toISOString();
    }

    this.projectService.updateProject(this.projectId, projectData).subscribe(
      response => {
        console.log('Project updated:', response);
        this.closeEditProjectModal();
      },
      error => {
        console.error('Error updating project:', error);
      }
    );
  }

  openEditProjectModal() {
    this.modalService.open(this, {
      ariaLabelledBy: 'editProjectModalLabel',
      backdrop: 'static',
      keyboard: false
    });
  }

  closeEditProjectModal(): void {
    this.modalService.dismissAll();
  }
}
