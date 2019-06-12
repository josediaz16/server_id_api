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

<style lang="scss">
.domain-item {
  background: #ffffff;
  border-radius: 15px;
  box-shadow: 1px 1px #f8fafb;
}

.domain-header {
  border-bottom: 2px solid #f8fafb;
  text-align: center;
  padding: 30px;

  img, object {
    width: 50px;
  }
}

.domain-detail,
.domain-servers {
  padding: 10px 30px;
}

.domain-detail {
  border-bottom: 2px solid #f8fafb;
  position: relative;

  span.down-arrow {
    box-shadow: 0px 0px 5px 1px #e8e1e1;
    border-radius: 50%;
    position: absolute;
    background: #fff;
    bottom: -10%;
    height: 40px;
    width: 40px;
    content: "";
    left: 45%;

    &:before {
      border-bottom: 2px solid #ccc;
      border-right: 2px solid #ccc;
      transform: rotate(45deg);
      display: inline-block;
      position: absolute;
      height: 8px;
      content: "";
      bottom: 39%;
      width: 8px;
      z-index: 2;
      left: 40%;
    }
  }

  li {
    align-items: center;
    display: flex;
    padding: 5px 0;
    font-size: 14px;
  }
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

  li {
    font-size: 12px;

    &.ip-address {
      text-overflow: ellipsis;
      white-space: nowrap;
      overflow: hidden;
      width: 150px;
    }
  }
}

.server-item {
  box-shadow: 1px 1px 5px #ccc;
  border-radius: 5px;
  padding: 5px 10px;
  width: 125px;
}

</style>
