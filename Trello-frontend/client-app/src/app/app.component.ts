import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { HttpClient, HttpClientModule } from '@angular/common/http';
import { RouterOutlet } from '@angular/router';
import { ProjectListComponent } from './components/project-list/project-list.component';
import { NavbarComponent } from './components/navbar/navbar.component';
import { LoginModalComponent } from './components/modals/login-modal/login-modal.component';
import { ProjectService } from './services/project.service';
import { Router } from '@angular/router';
import { FormsModule } from '@angular/forms'; 
import { RouterLink } from '@angular/router';
import { AuthService } from './services/auth.service';
import { RegisterModalComponent } from './components/modals/register-modal/register-modal.component';
import { VerifyModalComponent } from './components/modals/verify-modal/verify-modal.component';
import { EditProjectModalComponent } from './components/modals/edit-project-modal/edit-project-modal.component';
import { UserService } from './services/user.service';
import { TaskService } from './services/task.service';
import { NotificationService } from './services/notification.service';
import { ActivityService } from './services/activity.service';
import { AddProjectModalComponent } from './components/modals/add-project-modal/add-project-modal.component';
import { DashboardComponent } from './pages/dashboard/dashboard.component';
@Component({
  selector: 'app-root',
  standalone: true,
  imports: [
    CommonModule,
    HttpClientModule,
    RouterOutlet,
    ProjectListComponent,
    RouterLink,
    NavbarComponent,
    LoginModalComponent,
    RegisterModalComponent,
    VerifyModalComponent,
    EditProjectModalComponent,
    AddProjectModalComponent,
    DashboardComponent,
    FormsModule,  
    HttpClientModule
  ],
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
  providers: [UserService, ProjectService, AuthService,TaskService,NotificationService,ActivityService, Router, RouterLink],
})
export class AppComponent {
  title = 'Trello Clone';
}
