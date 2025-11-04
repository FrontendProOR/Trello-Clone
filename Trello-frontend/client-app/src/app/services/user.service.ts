import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { map, Observable } from 'rxjs';
import { User } from '../models/user.model';
import { AuthService } from './auth.service';
import { HttpHeaders } from '@angular/common/http';

@Injectable({
  providedIn: 'root',
})
export class UserService {
  private apiUrl = 'http://localhost:8080';

  constructor(private http: HttpClient, private authService: AuthService) {}

  getAllUsers(): Observable<User[]> {
    const token = this.authService.getToken();
    const headers = new HttpHeaders({
      Authorization: `Bearer ${token}`,
      'Content-Type': 'application/json',
    });

    const options = {
      headers: headers,
      withCredentials: true,
    };
    return this.http.get<User[]>(`${this.apiUrl}/users`, options);
  }

  getUserById(id: string): Observable<User> {
    const token = this.authService.getToken();
    const headers = new HttpHeaders({
      Authorization: `Bearer ${token}`,
      'Content-Type': 'application/json',
    });

    return this.http.get<User>(`${this.apiUrl}/users/${id}`, { headers }).pipe(
      map((user) => ({
        ...user,
        first_name: this.unescapeHTML(user.first_name),
        last_name: this.unescapeHTML(user.last_name),
        username: this.unescapeHTML(user.username),
        email: this.unescapeHTML(user.email),
      }))
    );
  }

  private unescapeHTML(html: string): string {
    const doc = new DOMParser().parseFromString(html, 'text/html');
    return doc.documentElement.textContent || '';
  }

  getUserByUsername(username: string): Observable<User> {
    const token = this.authService.getToken();
    const headers = new HttpHeaders({
      Authorization: `Bearer ${token}`,
      'Content-Type': 'application/json',
    });

    const options = {
      headers: headers,
      withCredentials: true,
    };
    return this.http.get<User>(
      `${this.apiUrl}/users/username/${username}`,
      options
    );
  }

  getUserByEmail(email: string): Observable<User> {
    const token = this.authService.getToken();
    const headers = new HttpHeaders({
      Authorization: `Bearer ${token}`,
      'Content-Type': 'application/json',
    });

    const options = {
      headers: headers,
      withCredentials: true,
    };

    return this.http.get<User>(`${this.apiUrl}/users/email/${email}`, options);
  }

  updateUser(id: string, user: User): Observable<any> {
    const token = this.authService.getToken();
    const headers = new HttpHeaders({
      Authorization: `Bearer ${token}`,
      'Content-Type': 'application/json',
    });

    const options = {
      headers: headers,
      withCredentials: true,
    };
    return this.http.put<User>(`${this.apiUrl}/users/${id}`, user, options);
  }

  deleteUser(id: string): Observable<void> {
    const token = this.authService.getToken();
    const headers = new HttpHeaders({
      Authorization: `Bearer ${token}`,
      'Content-Type': 'application/json',
    });

    const options = {
      headers: headers,
      withCredentials: true,
    };
    return this.http.delete<void>(`${this.apiUrl}/users/${id}`, options);
  }

  addUserToProject(projectId: string, userId: string): Observable<any> {
    const token = this.authService.getToken();
    const headers = new HttpHeaders({
      Authorization: `Bearer ${token}`,
      'Content-Type': 'application/json',
    });

    const options = {
      headers: headers,
      withCredentials: true,
    };
    return this.http.put<any>(
      `${this.apiUrl}/projects/${projectId}/users/${userId}`,
      {},
      options
    );
  }

  addUserToProjectByEmail(projectId: string, email: string): Observable<any> {
    const token = this.authService.getToken();
    const headers = new HttpHeaders({
      Authorization: `Bearer ${token}`,
      'Content-Type': 'application/json',
    });

    const options = {
      headers: headers,
      withCredentials: true,
    };
    return this.http.put<any>(
      `${this.apiUrl}/projects/${projectId}/users/${email}`,
      {},
      options
    );
  }

  removeUserFromProject(projectId: string, userId: string): Observable<any> {
    const token = this.authService.getToken();
    const headers = new HttpHeaders({
      Authorization: `Bearer ${token}`,
      'Content-Type': 'application/json',
    });

    const options = {
      headers: headers,
      withCredentials: true,
    };

    return this.http.delete<any>(
      `${this.apiUrl}/projects/${projectId}/users/${userId}`,
      options
    );
  }

  removeUserFromProjectByEmail(
    projectId: string,
    email: string
  ): Observable<any> {
    const token = this.authService.getToken();
    const headers = new HttpHeaders({
      Authorization: `Bearer ${token}`,
      'Content-Type': 'application/json',
    });

    const options = {
      headers: headers,
      withCredentials: true,
    };

    return this.http.delete<any>(
      `${this.apiUrl}/projects/${projectId}/users/${email}`,
      options
    );
  }

  getAllMembersInProject(projectId: string): Observable<User[]> {
    const token = this.authService.getToken();
    const headers = new HttpHeaders({
      Authorization: `Bearer ${token}`,
      'Content-Type': 'application/json',
    });

    const options = {
      headers: headers,
      withCredentials: true,
    };

    return this.http.get<User[]>(
      `${this.apiUrl}/projects/${projectId}/members`,
      options
    );
  }
}
