import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Notification } from '../models/notification.model';
import { AuthService } from './auth.service';

@Injectable({
  providedIn: 'root'
})
export class NotificationService {
  private apiUrl = 'http://localhost:8084';

  constructor(private http: HttpClient, private authService: AuthService) { }

  // Dobijanje svih notifikacija
  getNotifications(userId: string): Observable<Notification[]> {
    const token = this.authService.getToken();
    console.log('Authorization token:', token); // Debug

    const headers = new HttpHeaders({
      'Authorization': `Bearer ${token}`
    });

    return this.http.get<Notification[]>(`${this.apiUrl}/notifications/${userId}`, { headers });
  }

  // Označavanje notifikacije kao pročitane
  markNotificationAsRead(notificationId: string): Observable<void> {
    const token = this.authService.getToken();
    const headers = new HttpHeaders({
      'Authorization': `Bearer ${token}`
    });

    return this.http.put<void>(`${this.apiUrl}/notifications/read`, { notificationId }, { headers });
  }
}
