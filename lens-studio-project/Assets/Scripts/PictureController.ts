import { SIK } from "SpectaclesInteractionKit.lspkg/SIK";

@component
export class PictureController extends BaseScriptComponent {
  @input scannerPrefab: ObjectPrefab;

  private isEditor = global.deviceInfoSystem.isEditor();

  private rightHand = SIK.HandInputData.getHand("right");
  private leftHand = SIK.HandInputData.getHand("left");

  private leftDown = false;
  private rightDown = false;

  onAwake() {
    this.rightHand.onPinchUp.add(this.rightPinchUp);
    this.rightHand.onPinchDown.add(this.rightPinchDown);
    this.leftHand.onPinchUp.add(this.leftPinchUp);
    this.leftHand.onPinchDown.add(this.leftPinchDown);
    if (this.isEditor) {
      this.createEvent("TouchStartEvent").bind(this.editorTest.bind(this));
    } else {
      var obj = this.getSceneObject();
      if (obj.getChildrenCount() > 0) {
        obj.getChild(0).destroy();
      }
    }
  }

  editorTest() {
    print("Creating Editor Scanner...");
    this.createScanner();
  }

  private leftPinchDown = () => {
    this.leftDown = true;
    if (this.rightDown && this.isPinchClose()) {
      this.createScanner();
    }
  };

  private leftPinchUp = () => {
    this.leftDown = false;
  };

  private rightPinchDown = () => {
    this.rightDown = true;
    if (this.leftDown && this.isPinchClose()) {
      this.createScanner();
    }
  };

  private rightPinchUp = () => {
    this.rightDown = false;
  };

  isPinchClose() {
    return this.leftHand.thumbTip.position.distance(this.rightHand.thumbTip.position) < 10;
  }

  createScanner() {
    var scanner = this.scannerPrefab.instantiate(this.getSceneObject());
  }
}
