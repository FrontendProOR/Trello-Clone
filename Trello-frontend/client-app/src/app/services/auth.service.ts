import { Injectable, PLATFORM_ID } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { map, Observable } from 'rxjs';
import { isPlatformBrowser } from '@angular/common';
import { Inject } from '@angular/core';
import { HttpHeaders } from '@angular/common/http';

@Injectable({
  providedIn: 'root',
})
export class AuthService {
  private apiUrl = 'http://localhost:8080/users';

  constructor(
    @Inject(PLATFORM_ID) private platformId: Object,
    private http: HttpClient
  ) {}

  login(email: string, password: string): Observable<any> {
    return this.http.post<any>(this.apiUrl + '/login', { email, password });
  }

  checkPasswordBlacklist(password: string): Observable<boolean> {
    return this.http
      .post<{ blacklisted: boolean }>(`${this.apiUrl}/check-password`, {
        password,
      })
      .pipe(
        map((response) => response.blacklisted) // Izvlači `blacklisted` vrednost
      );
  }

  register(user: any): Observable<any> {
    return this.http.post<any>(this.apiUrl + '/register', user);
  }

  verify(email: string, code: string): Observable<any> {
    return this.http.post<any>(this.apiUrl + '/verify', { email, code });
  }

  setToken(token: string): void {
    localStorage.setItem('token', token);
  }

  getToken() {
    if (typeof window !== 'undefined' && window.localStorage) {
      return localStorage.getItem('token');
    }
    return null;
  }

  changePassword(
    email: string,
    oldPassword: string,
    newPassword: string
  ): Observable<any> {
    const token = this.getToken();
    const headers = new HttpHeaders({
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    });

    const options = {
      headers: headers,
      withCredentials: true,
    };

    return this.http.put<any>(
      `${this.apiUrl}/change-password`,
      { email, oldPassword, newPassword },
      options
    );
  }

  removeToken(): void {
    localStorage.removeItem('token');
  }

  isLoggedIn(): boolean {
    if (isPlatformBrowser(this.platformId)) {
      return !!localStorage.getItem('token');
    } else {
      return false;
    }
  }

  logout(): void {
    this.removeToken();
  }

  sendMagicLink(email: string) {
    return this.http.post(`${this.apiUrl}/magic-link`, { email });
  }

  magicLogin(token: string) {
    return this.http.get(`${this.apiUrl}/magic-link/validate?token=${token}`);
  }
  generateCaptcha(): Observable<{ id: string; image: string }> {
    return this.http.get<{ id: string; image: string }>(
      `${this.apiUrl}/captcha`
    );
  }

  registerWithCaptcha(data: any): Observable<any> {
    return this.http.post(`${this.apiUrl}/register`, data);
  }
  validateCaptcha(captchaId: string, captchaText: string): Observable<any> {
    return this.http.post(`${this.apiUrl}/captcha/validate`, { captchaId, captchaText });
  }
  
}
