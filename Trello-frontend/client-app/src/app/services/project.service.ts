import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { map, Observable } from 'rxjs';
import { AuthService } from './auth.service';

@Injectable({
  providedIn: 'root',
})
export class ProjectService {
  private apiUrl = 'http://localhost:8080/projects';

  constructor(private http: HttpClient, private authService: AuthService) {}

  getProjects(): Observable<any> {
    return this.http.get<any>(this.apiUrl);
  }

  getProjectById(id: string): Observable<any> {
    const token = this.authService.getToken();
    const headers = new HttpHeaders({
      Authorization: `Bearer ${token}`,
      'Content-Type': 'application/json',
    });
    return this.http.get<any>(`${this.apiUrl}/${id}`, { headers });
  }

  addProject(project: any): Observable<any> {
    const token = this.authService.getToken();
    const headers = new HttpHeaders({
      Authorization: `Bearer ${token}`,
      'Content-Type': 'application/json',
    });

    const options = {
      headers: headers,
      withCredentials: true,
    };

    return this.http.post<any>(`${this.apiUrl}`, project, options);
  }

  updateProject(id: string, project: any): Observable<any> {
    const token = this.authService.getToken();
    const headers = new HttpHeaders({
      Authorization: `Bearer ${token}`,
      'Content-Type': 'application/json',
    });

    const options = {
      headers: headers,
      withCredentials: true,
    };
    return this.http.put<any>(`${this.apiUrl}/${id}`, project, options);
  }

  deleteProject(projectId: string): Observable<any> {
    const token = this.authService.getToken();
    const headers = new HttpHeaders({
      Authorization: `Bearer ${token}`,
      'Content-Type': 'application/json',
    });
    return this.http.delete<any>(`${this.apiUrl}/${projectId}`, { headers });
  }

  getProjectsByMemberId(userId: string): Observable<any[]> {
    const token = localStorage.getItem('token');
    const headers = new HttpHeaders({
      Authorization: `Bearer ${token}`,
      'Content-Type': 'application/json',
    });
    return this.http
      .get<any[]>(`${this.apiUrl}/member/${userId}`, { headers })
      .pipe(
        map((projects) => {
          return projects.map((project) => {
            project.description = this.unescapeHTML(project.description);
            project.name = this.unescapeHTML(project.name);

            return project;
          });
        })
      );
  }

  getProjectsByManagerId(userId: string): Observable<any[]> {
    const token = localStorage.getItem('token');
    const headers = new HttpHeaders({
      Authorization: `Bearer ${token}`,
      'Content-Type': 'application/json',
    });
    return this.http
      .get<any[]>(`${this.apiUrl}/manager/${userId}`, {
        headers,
      })
      .pipe(
        map((projects) => {
          return projects.map((project) => {
            project.description = this.unescapeHTML(project.description);
            project.name = this.unescapeHTML(project.name);
            return project;
          });
        })
      );
  }
  private unescapeHTML(html: string): string {
    const doc = new DOMParser().parseFromString(html, 'text/html');
    return doc.documentElement.textContent || '';
  }
}
