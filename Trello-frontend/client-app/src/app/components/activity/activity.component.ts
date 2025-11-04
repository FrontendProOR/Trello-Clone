import { Component } from '@angular/core';
import { ActivityService } from '../../services/activity.service';
@Component({
  selector: 'app-activity',
  standalone: true,
  imports: [],
  templateUrl: './activity.component.html',
  styleUrl: './activity.component.scss'
})
export class ActivityComponent {
  activities: any[] = [];
  
  constructor(private activityService: ActivityService) { }
  
  ngOnInit(): void {
    const userId = localStorage.getItem('userId') || '';

    this.activityService.getActivities(userId).subscribe(data => {
      this.activities = data;
    });
  }
}
