import { ComponentFixture, TestBed } from '@angular/core/testing';

import { AppForgotPasswordModalComponent } from './app-forgot-password-modal.component';

describe('AppForgotPasswordModalComponent', () => {
  let component: AppForgotPasswordModalComponent;
  let fixture: ComponentFixture<AppForgotPasswordModalComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AppForgotPasswordModalComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(AppForgotPasswordModalComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
