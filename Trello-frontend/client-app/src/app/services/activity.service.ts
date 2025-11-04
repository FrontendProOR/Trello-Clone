import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})

export class ActivityService {
  private apiUrl = 'http://localhost:8080/api';

  constructor( private http: HttpClient) { }

  getActivities(userId: string): Observable<any> {
    return this.http.get(`${this.apiUrl}/activities/${userId}`);
  }

}
