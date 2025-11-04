import { ComponentFixture, TestBed } from '@angular/core/testing';

import { RemoveMemberModalComponent } from './remove-member-modal.component';

describe('RemoveMemberModalComponent', () => {
  let component: RemoveMemberModalComponent;
  let fixture: ComponentFixture<RemoveMemberModalComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [RemoveMemberModalComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(RemoveMemberModalComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
