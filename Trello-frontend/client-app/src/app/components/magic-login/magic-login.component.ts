import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { AuthService } from '../../services/auth.service';  // Assuming you have an AuthService

@Component({
  selector: 'app-magic-login',
  templateUrl: './magic-login.component.html',
  styleUrls: ['./magic-login.component.scss'],
})
export class MagicLoginComponent implements OnInit {
  constructor(
    private route: ActivatedRoute,
    private authService: AuthService,  // AuthService to handle backend calls
    private router: Router
  ) {}

  ngOnInit(): void {
    // Get token from query parameters
    this.route.queryParams.subscribe((params) => {
      const token = params['token'];

      if (token) {
        this.authService.magicLogin(token).subscribe(
          (response: any) => {
            // Store JWT token in localStorage or sessionStorage
            localStorage.setItem('token', response.token);
            // Redirect to dashboard
            this.router.navigate(['/dashboard']);
          },
          (error: any) => {
            console.error('Magic login failed', error);
            this.router.navigate(['/login']);
          }
        );
      } else {
        this.router.navigate(['/login']);  // No token provided
      }
    });
  }
}
