import { definePreset } from '@primeuix/themes';
import Aura from '@primeuix/themes/aura';

const MyPreset = definePreset(Aura, {
  semantic: {
    primary: {
      50:  '#E6F0F7', // even less saturated very light blue
      100: '#C6D9EB',
      200: '#A7C2DF',
      300: '#88ABD3',
      400: '#6A96C7',
      500: '#4F82BB',
      600: '#436FA3', // less saturated base blue
      700: '#385D8B',
      800: '#2C4A6F',
      900: '#1C3047',
      950: '#101A24'  // even deeper, less saturated dark blue
    }
  }
});






export default {
    preset: MyPreset,
    options: {
        darkModeSelector: '.p-dark'
    }
};
