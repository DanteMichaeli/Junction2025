@component
export class MoveUpDown extends BaseScriptComponent {
    // 5 cm up and 5 cm down
    @input
    amplitude: number = 5.0;

    // Full cycle (up+down) in seconds
    @input
    period: number = 1.5;

    private startPos: vec3 | null = null;
    private startTime: number = 0;

    onAwake() {
        const tr = this.sceneObject.getTransform();
        if (!tr) {
            print("MoveUpDown: No Transform on this SceneObject");
            return;
        }

        this.startPos = tr.getWorldPosition();
        this.startTime = getTime();

        // Register update event (per docs)
        this.createEvent('UpdateEvent').bind(this.onUpdate.bind(this));
    }

    private onUpdate() {
        if (!this.startPos) {
            return;
        }

        const tr = this.sceneObject.getTransform();
        if (!tr) {
            return;
        }

        // Elapsed time in seconds
        const elapsed = getTime() - this.startTime;

        // Map time → phase (0..2π over "period" seconds)
        const phase = (elapsed / this.period) * Math.PI * 2.0;

        // Sine wave between -amplitude and +amplitude
        const offset = Math.sin(phase) * this.amplitude;

        const newPos = new vec3(
            this.startPos.x,
            this.startPos.y + offset, // move up/down in cm
            this.startPos.z
        );

        tr.setWorldPosition(newPos);
    }
}