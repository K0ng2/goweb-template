import { library, config } from '@fortawesome/fontawesome-svg-core'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import {
} from "@fortawesome/free-solid-svg-icons";

import {
} from "@fortawesome/free-brands-svg-icons";

const icons = [
];


// This is important, we are going to let Nuxt worry about the CSS
config.autoAddCss = false;

// You can add your icons directly in this plugin. See other examples for how you
// can add other styles or just individual icons.
library.add(...icons);

export default defineNuxtPlugin((nuxtApp) => {
	nuxtApp.vueApp.component('FontAwesomeIcon', FontAwesomeIcon)
})
