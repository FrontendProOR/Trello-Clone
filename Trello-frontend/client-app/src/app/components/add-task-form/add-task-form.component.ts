import { Component } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule} from '@angular/forms';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-add-task-form',
  standalone: true, 
  imports: [CommonModule, ReactiveFormsModule], 
  templateUrl: './add-task-form.component.html',
  styleUrls: ['./add-task-form.component.scss'],
})

export class AddTaskComponent {
  taskForm: FormGroup;
  projectId: string | null = null;

  constructor(
    private fb: FormBuilder,
    private route: ActivatedRoute,
    private router: Router
  ) {
    this.projectId = this.route.snapshot.paramMap.get('id');
    this.taskForm = this.fb.group({
      name: ['', Validators.required],
      description: [''],
      assignedTo: ['672ab1afcf88771952544e51'],
      dueDate: ['', Validators.required],
      status: ['Pending', Validators.required],
      priority: ['Medium', Validators.required],
    });
  }

  onSubmit(): void {
    if (this.taskForm.valid && this.projectId) {
      const taskData = { ...this.taskForm.value, projectId: this.projectId };
      console.log('Task Data to be sent:', taskData);
      this.router.navigate([`/projects/${this.projectId}`]);
    }
  }

}
