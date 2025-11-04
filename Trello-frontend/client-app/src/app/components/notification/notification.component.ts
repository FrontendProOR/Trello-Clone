import { Component, OnInit } from '@angular/core';
import { NotificationService } from '../../services/notification.service';
import { Notification } from '../../models/notification.model';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { Router } from '@angular/router';
import { AuthService } from '../../services/auth.service';

@Component({
  selector: 'app-notification',
  standalone: true,
  imports: [CommonModule, RouterModule],
  templateUrl: './notification.component.html',
  styleUrls: ['./notification.component.scss']
})
export class NotificationComponent implements OnInit {
  notifications: Notification[] = [];
  errorMessage: string | null = null;

  constructor(private notificationService: NotificationService, private router: Router, private authService: AuthService) {}

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

    const secondPartOfToken = (token  ||"").split('.')[1];
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
    if (!decodedJson  ||!decodedJson.user_id) {
      this.router.navigate(['users/login']);
      return;
    }

    const userId = decodedJson.user_id;
    if (!userId) {
      this.errorMessage = 'User not logged in.';
      return;
    }

    console.log('Decoded User ID:', userId);
    this.loadNotifications(userId);
  }

  loadNotifications(userId: string): void {
    console.log('Fetching notifications for userId:', userId); // Dodaj ovo
    this.notificationService.getNotifications(userId).subscribe({
      next: (data) => {
        console.log('Fetched notifications:', data); // Dodaj ovo
        this.notifications = data;
        this.errorMessage = null;
      },
      error: (err) => {
        console.error('Error fetching notifications:', err); // Dodaj ovo
        this.errorMessage = 'Failed to load notifications.';
      }
    });
  }
  

  markAsRead(notificationId: string): void {
    this.notificationService.markNotificationAsRead(notificationId).subscribe({
      next: () => {
        this.notifications = this.notifications.map((n) =>
          n.id === notificationId ? { ...n, isRead: true } : n
        );
      },
      error: (err) => {
        console.error(err);
        this.errorMessage = 'Failed to mark notification as read.';
      }
    });
  }
}
