<script lang="ts">

import { notify } from '@viamrobotics/prime';
import { displayError } from '@/lib/error';
import { rcLogConditionally } from '@/lib/log';
import { BoardClient, type ServiceError } from '@viamrobotics/sdk';
import Collapse from '@/lib/components/collapse.svelte';
import { useRobotClient } from '@/hooks/robot-client';

export let name: string;
export let status: undefined | {
  analogs: Record<string, { value: number }>
  digital_interrupts: Record<string, { value: number }>
};

const { robotClient } = useRobotClient();

const boardClient = new BoardClient($robotClient, name, { requestLogger: rcLogConditionally });

let getPin = '';
let setPin = '';
let setLevel = '';
let pwm = '';
let pwmFrequency = '';
let getPinMessage = '';
let analogPin = '';
let analogValue = '';

const getGPIO = async () => {
  try {
    const isHigh = await boardClient.getGPIO(getPin);
    getPinMessage = `Pin: ${getPin} is ${isHigh ? 'high' : 'low'}`;
  } catch (error) {
    notify.danger((error as ServiceError).message);
  }
};

const setGPIO = async () => {
  try {
    await boardClient.setGPIO(setPin, setLevel === 'high');
  } catch (error) {
    displayError(error as ServiceError);
  }
};

const getPWM = async () => {
  try {
    const dutyCyclePct = await boardClient.getPWM(getPin);
    getPinMessage = `Pin ${getPin}'s duty cycle is ${dutyCyclePct * 100}%.`;
  } catch (error) {
    displayError(error as ServiceError);
  }
};

const setPWM = async () => {
  try {
    await boardClient.setPWM(setPin, Number.parseFloat(pwm) / 100);
  } catch (error) {
    displayError(error as ServiceError);
  }
};

const getPWMFrequency = async () => {
  try {
    const frequencyHz = await boardClient.getPWMFrequency(getPin);
    getPinMessage = `Pin ${getPin}'s frequency is ${frequencyHz}Hz.`;
  } catch (error) {
    displayError(error as ServiceError);
  }
};

const setPWMFrequency = async () => {
  try {
    await boardClient.setPWMFrequency(setPin, Number.parseFloat(pwmFrequency));
  } catch (error) {
    displayError(error as ServiceError);
  }
};

const writeAnalog = async () => {
  try {
    await boardClient.writeAnalog(analogPin, Number.parseInt(analogValue, 10));
  } catch (error) {
    displayError(error as ServiceError);
  }
};

const handleGetPinInput = (event: CustomEvent<{ value: string}>) => {
  getPin = event.detail.value;
};

const handleSetPinInput = (event: CustomEvent<{ value: string}>) => {
  setPin = event.detail.value;
};

const handlePwmInput = (event: CustomEvent<{ value: string}>) => {
  pwm = event.detail.value;
};

const handlePwmFrequencyInput = (event: CustomEvent<{ value: string}>) => {
  pwmFrequency = event.detail.value;
};

const handleAnalogPinInput = (event:CustomEvent<{ value: string}>) => {
  analogPin = event.detail.value;
};

const handleAnalogValueInput = (event:CustomEvent<{ value: string}>) => {
  analogValue = event.detail.value;
};

</script>

<Collapse title={name}>
  <v-breadcrumbs
    slot="title"
    crumbs="board"
  />
  <div class="overflow-auto border border-t-0 border-medium p-4">
    <h3 class="mb-2">
      Analog Readers
    </h3>
    <table class="mb-4 table-auto border border-medium">
      {#each Object.entries(status?.analogs ?? {}) as [analogName, analog] (analogName)}
        <tr>
          <th class="border border-medium p-2">
            {analogName}
          </th>
          <td class="border border-medium p-2">
            {analog.value || 0}
          </td>
        </tr>
      {/each}
    </table>

    <h3 class="mb-2">
      Digital Interrupts
    </h3>
    <table class="mb-4 w-full table-auto border border-medium">
      {#each Object.entries(status?.digital_interrupts ?? {}) as [interruptName, interrupt] (interruptName)}
        <tr>
          <th class="border border-medium p-2">
            {interruptName}
          </th>
          <td class="border border-medium p-2">
            {interrupt.value || 0}
          </td>
        </tr>
      {/each}
    </table>

    <h3 class="mb-2">
      GPIO
    </h3>
    <table class="mb-4 w-full table-auto border border-medium">
      <tr>
        <th class="border border-medium p-2">
          Get
        </th>
        <td class="border border-medium p-2">
          <div class="flex flex-wrap items-end gap-2">
            <v-input
              label="Pin"
              type="text"
              value={getPin}
              on:input={handleGetPinInput}
            />
            <v-button
              label="Get Pin State"
              on:click={getGPIO}
            />
            <v-button
              label="Get PWM Duty Cycle"
              on:click={getPWM}
            />
            <v-button
              label="Get PWM Frequency"
              on:click={getPWMFrequency}
            />
            <span class="py-2">
              {getPinMessage}
            </span>
          </div>
        </td>
      </tr>

      <tr>
        <th class="border border-medium p-2">
          Set
        </th>
        <td class="p-2">
          <div class="flex flex-wrap items-end gap-2">
            <v-input
              value={setPin}
              type="text"
              class="mr-2"
              label="Pin"
              on:input={handleSetPinInput}
            />
            <select
              bind:value={setLevel}
              class="mr-2 h-[30px] border border-medium bg-white text-sm"
            >
              <option>low</option>
              <option>high</option>
            </select>
            <v-button
              class="mr-2"
              label="Set Pin State"
              on:click={setGPIO}
            />
            <v-input
              value={pwm}
              label="PWM Duty Cycle"
              type="number"
              class="mr-2"
              on:input={handlePwmInput}
            />
            <v-button
              class="mr-2"
              label="Set PWM Duty Cycle"
              on:click={setPWM}
            />
            <v-input
              value={pwmFrequency}
              label="PWM Frequency"
              type="number"
              class="mr-2"
              on:input={handlePwmFrequencyInput}
            />
            <v-button
              class="mr-2"
              label="Set PWM Frequency"
              on:click={setPWMFrequency}
            />
          </div>
        </td>
      </tr>
    </table>
    <h3 class="mb-2">
      Analog Write
     </h3>
     <table class="mb-4 w-full table-auto border border-medium">
     <td class="border border-medium p-2">
       <div class="flex flex-wrap items-end gap-2">
         <v-input
           label="Pin"
           type="text"
           value={getPin}
           on:input={handleAnalogPinInput}
         />
         <v-input
         label="Value"
         type="text"
         value={getPin}
         on:input={handleAnalogValueInput}
       />
       <v-button
       class="mr-2"
       label="Set Analog Value"
       on:click={writeAnalog}
     />
       </div>
     </td>
     </table>
  </div>
</Collapse>
