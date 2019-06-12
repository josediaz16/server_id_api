<template>
  <svg class="icon">
    <use :xlink:href="iconPath"/>
  </svg>
</template>

<script>
const icons = require("../assets/icons.svg")

export default {
  name: "Icon",
  props: {
    iconname: String,
  },
  computed: {
    iconPath: function() {
      return `${icons}#${this.iconname}`
    }
  }
}
</script>
