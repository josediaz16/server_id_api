<template>

  <div class="domain-item">
    <div class="domain-header">
      <img v-bind:src="logo" @error="onLogoFailure($event)">
      <h4>{{name}}</h4>
    </div>

    <div class="domain-detail">
      <h5>Checklist</h5>
      <ul>

        <li>
          <Icon v-bind:iconname="checkOrNot(!domain.is_down)"/>
          Domain Up
        </li>

        <li>
          <Icon v-bind:iconname="checkOrNot(!domain.servers_changed)"/>
          Servers changed
        </li>

        <li>
          <Icon v-bind:iconname="sslGradeIcon(domain.ssl_grade)"/>
          Ssl grade {{domain.ssl_grade}}
        </li>

        <li>
          <Icon v-bind:iconname="sslGradeIcon(domain.previous_ssl_grade)"/>
          Previous Ssl grade {{domain.ssl_grade}}
        </li>
      </ul>
      <span v-on:click="showServers = !showServers" class="down-arrow"></span>
    </div>


    <div v-if="showServers" class="domain-servers">
      <h5>Servers</h5>
      <div class="server-list">
        <div v-for="(server) in domain.servers" v-bind:key="server.ip_address" class="server-item">
          <ul>
            <li class="ip-address">
              <Icon iconname="placeholder"/>
              {{server.ip_address}}
            </li>
            <li>
              <Icon v-bind:iconname="sslGradeIcon(server.grade)"/>
              Ssl grade {{server.grade}}
            </li>
            <li>
              <Icon iconname="flag"/>
              Country: {{server.country}}
            </li>
            <li>
              <Icon iconname="visitor"/>
              Owner: {{server.owner}}
            </li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</template>

<script>

import Icon from './Icon'

const icons = require("../assets/icons.svg")
const DefaultLogo = require("../assets/default-logo.png")

const GradeOptions = {
  "A":  "A",
  "A-": "A",
  "A+": "A",
  "B":  "B",
  "C":  "C",
  "F":  "F",
  "U":  "U",
  "":   "U"
}

export default {
  name: "DomainCard",
  props: {
    domain: Object,
    name: String
  },
  components: {
    Icon
  },
  data: function() {
    return {
      showServers: false
    }
  },
  methods: {
    checkOrNot: function(flag) {
      return flag ? "checked" : "error"
    },
    sslGradeIcon: function(grade) {
      return `letter-${GradeOptions[grade].toLowerCase()}`
    },
    onLogoFailure(event) {
      event.target.src = DefaultLogo
    }
  },
  computed: {
    logo: function() {
      return this.domain.logo === "" ? DefaultLogo : this.domain.logo
    },
  }
}
</script>

<style>
.domain-item {
  background: #ffffff;
  border-radius: 15px;
  box-shadow: 1px 1px #f8fafb;
}

.domain-header {
  border-bottom: 2px solid #f8fafb;
  text-align: center;
  padding: 30px;
}

.domain-header img,
.domain-header object {
  width: 50px;
}

.domain-detail,
.domain-servers {
  padding: 10px 30px;
}

.domain-detail {
  border-bottom: 2px solid #f8fafb;
  position: relative;
}

.domain-detail span.down-arrow {
  content: "";
  width: 40px;
  height: 40px;
  position: absolute;
  bottom: -10%;
  left: 45%;
  border-radius: 50%;
  background: #fff;
  box-shadow: 0px 0px 5px 1px #e8e1e1;
}

.domain-detail span.down-arrow:before {
  content: "";
  transform: rotate(45deg);
  width: 8px;
  height: 8px;
  display: inline-block;
  position: absolute;
  border-bottom: 2px solid #ccc;
  border-right: 2px solid #ccc;
  bottom: 39%;
  left: 40%;
  z-index: 2;
}

.domain-detail li {
  align-items: center;
  display: flex;
  padding: 5px 0;
  font-size: 14px;
}

.domain-detail svg.icon,
.domain-servers svg.icon {
  height: 20px;
  width: 20px;
  margin-right: 5px;
}

.server-list {
  grid-template-columns: 1fr 1fr;
  grid-gap: 10px;
  display: grid;
}

.server-list li {
  font-size: 12px;
}

.server-list li.ip-address {
  width: 150px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.server-item {
  box-shadow: 1px 1px 5px #ccc;
  border-radius: 5px;
  padding: 5px 10px;
  width: 125px;
}

</style>
