import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { HomepageComponent } from './pages/homepage/homepage.component';
import { DashboardComponent } from './pages/dashboard/dashboard.component';
import { ProjectPageComponent } from './pages/project-page/project-page.component';
import { ProfilePageComponent } from './pages/profile-page/profile-page.component';
import { TaskPageComponent } from './pages/task-page/task-page.component';
import { LoginModalComponent } from './components/modals/login-modal/login-modal.component';
import { RegisterModalComponent } from './components/modals/register-modal/register-modal.component';
import { VerifyModalComponent } from './components/modals/verify-modal/verify-modal.component';
import { AddProjectModalComponent } from './components/modals/add-project-modal/add-project-modal.component';
import { EditProjectModalComponent } from './components/modals/edit-project-modal/edit-project-modal.component';
import { ActivityComponent } from './components/activity/activity.component';
import { NotificationComponent } from './components/notification/notification.component';
import { ProjectListComponent } from './components/project-list/project-list.component';
import { MagicLoginComponent } from './components/magic-login/magic-login.component';


export const routes: Routes = [
  { path: '', component: HomepageComponent },
  { path: 'dashboard',component: DashboardComponent},
  { path: 'projects/manager/:id', component: ProjectListComponent},
  { path: 'projects/member/:id', component: ProjectListComponent},
  { path: 'users/profile', component: ProfilePageComponent},
  { path: 'users/profile/:id', component: ProfilePageComponent},
  { path: 'projects/:id/tasks/:taskId', component: TaskPageComponent},
  { path: 'notifications', component: NotificationComponent },
  { path: 'projects/:id/activity', component: ActivityComponent },
  { path: 'users/login', component: LoginModalComponent},
  { path: 'users/register', component: RegisterModalComponent},
  { path: 'users/verify', component: VerifyModalComponent},
  // { path: 'projects', component: AddProjectModalComponent },
  { path: 'projects/:id', component: EditProjectModalComponent },
  { path: 'projects/:id/tasks', component: ProjectPageComponent },
  { path: 'users/magic-login', component: MagicLoginComponent },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule],
})

export class AppRoutingModule {}
