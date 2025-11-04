import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { TaskService } from '../../services/task.service';  
import { Task } from '../../models/task.model';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormsModule } from '@angular/forms';

@Component({
  selector: 'app-project-page',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, FormsModule],
  templateUrl: './project-page.component.html',
  styleUrls: ['./project-page.component.scss']
})
export class ProjectPageComponent implements OnInit {
  projectId: string = '';
  tasks: Task[] = [];

  constructor(
    private taskService: TaskService,
    private route: ActivatedRoute
  ) {}

  ngOnInit(): void {
    this.route.paramMap.subscribe(params => {
      this.projectId = params.get('id')!;
      this.loadProjectTasks();
    });
  }
  
  loadProjectTasks() {
    this.taskService.getTasksByProjectId(this.projectId).subscribe(tasks => {
      console.log('Fetched tasks:', tasks); 
      this.tasks = tasks;
    }, error => {
      console.error('Error fetching tasks:', error);
    });
  }

  getStatusClass(status: string): string {
    switch (status) {
      case 'Completed':
        return 'completed';
      case 'Pending':
        return 'pending';
      case 'In Progress':
        return 'inProgress';
      default:
        return '';
    }
  }
  
  getPriorityClass(priority: string): string {
    switch (priority) {
      case 'High':
        return 'highPriority';
      case 'Medium':
        return 'mediumPriority';
      case 'Low':
        return 'lowPriority';
      default:
        return '';
    }
  }

  deleteTask(taskId: string) {
    if (confirm('Da li ste sigurni da želite da obrišete ovaj task?')) {
      this.taskService.deleteTask(taskId).subscribe(
        () => {
          this.tasks = this.tasks.filter(task => task.id !== taskId); // Promenili _id u id
          alert('Task je uspešno obrisan.');
        },
        (error: any) => {
          console.error('Greška pri brisanju taska:', error);
          alert('Došlo je do greške pri brisanju taska.');
        }
      );
    }
  }
}
