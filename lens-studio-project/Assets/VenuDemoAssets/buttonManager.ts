import { CapsuleButton } from "SpectaclesUIKit.lspkg/Scripts/Components/Button/CapsuleButton";
import { LineSwitcher } from "./lineSwitcher";

@component
export class LineSelectionManager extends BaseScriptComponent {

    @input buttonR: CapsuleButton
    @input buttonW: CapsuleButton
    @input buttonC: CapsuleButton

    @input switcher: LineSwitcher

    private mask = 0;   // stores which buttons are active

    onAwake() {
        this.buttonR.onTriggerDown.add(() => {
            this.toggleBit(1);   // R = 1
        });

        this.buttonW.onTriggerDown.add(() => {
            this.toggleBit(2);   // W = 2
        });

        this.buttonC.onTriggerDown.add(() => {
            this.toggleBit(4);   // C = 4
        });
    }

    toggleBit(bit: number) {
        // Toggle the bit using XOR
        this.mask ^= bit;

        print(`This is the current bitmask: ${this.mask}`);

        // Update the lines
        this.switcher.showCombination(this.mask);
    }
}
