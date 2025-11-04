import { Component } from '@angular/core';
import { ProjectListComponent } from '../../components/project-list/project-list.component';
import { FormsModule } from '@angular/forms';
import { AddProjectModalComponent } from '../../components/modals/add-project-modal/add-project-modal.component';

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [ProjectListComponent, FormsModule, AddProjectModalComponent],
  templateUrl: './dashboard.component.html',
  styleUrl: './dashboard.component.scss'
})
export class DashboardComponent {

}
