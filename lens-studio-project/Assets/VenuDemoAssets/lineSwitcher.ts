@component
export class LineSwitcher extends BaseScriptComponent {

    @input
    lines: SceneObject[];   

    onAwake() {
        this.hideAll();   
    }

    showLine(index: number) {
        if (!this.lines) return;

        for (let i = 0; i < this.lines.length; i++) {
            const line = this.lines[i];
            if (line) {
                line.enabled = (i === index);
            }
        }
    }

    showCombination(mask: number) {
        if (!this.lines) return;

        this.hideAll();

        if (mask === 0) return;  // nothing selected

        let index = mask - 1

        print(`This is the current index: ${index}`);
        
        if (this.lines[index]) {
            this.lines[index].enabled = true;
        }
    }

    hideAll() {
        if (!this.lines) return;

        for (let i = 0; i < this.lines.length; i++) {
            const line = this.lines[i];
            if (line) {
                line.enabled = false;
            }
        }
    }
}
