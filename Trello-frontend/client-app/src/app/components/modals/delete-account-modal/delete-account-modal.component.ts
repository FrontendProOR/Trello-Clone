import { Component } from '@angular/core';
import { NgbActiveModal } from '@ng-bootstrap/ng-bootstrap';

@Component({
  selector: 'app-delete-account-modal',
  templateUrl: './delete-account-modal.component.html',
  styleUrls: ['./delete-account-modal.component.scss']
})
export class DeleteAccountModalComponent {
  constructor(public activeModal: NgbActiveModal) {}

  // This will close the modal without doing anything (cancel)
  cancel() {
    this.activeModal.dismiss();
  }

  // This will confirm the deletion
  confirmDeletion() {
    this.activeModal.close('confirmed');
  }
}
