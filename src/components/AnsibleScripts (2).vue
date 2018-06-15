<template>
  <div class="row h-100 w-100">
    <div class="col-md-4 scroll h-100">
      <table class="table table-hover table-bordered table-sm table-dark">
        <tbody>
        <tr v-for="script in ansiblescripts">
          <th scope="row" style="width: 5%">{{ script.id }}</th>
          <!--<td>{{ script.name }}</td> -->
          <td style="width: 5%">{{ script.path }}</td>
          <!-- <td>{{ script.description }}</td> -->
          <td style="width: 5%">
            <button type="button" class="btn btn-sm  btn-outline-success" v-on:click="runAnsibleScript(script)">
                <fa-icon icon="play" aria-hidden="true"/>
            </button>

          </td>
        </tr>
        </tbody>
      </table>
    </div>
    <div class="col-md-8 shell h-100" ref="shell">
      <i v-show="loading" class="fa fa-spinner fa-spin fa-4x"></i>
      <span class="h-100" v-html="script_result "></span>
    </div>
  </div>
</template>

<script>
export default {
  name: 'AnsibleScripts',
  data: function () {
    return {
      script_result: '',
      ansiblescripts: [],
      loading: false
    }
  },
  mounted: function () {
    this.$http.get('/ansiblescripts').then(function (response) {
      this.ansiblescripts = response.data.items ? response.data.items : []
    })
    console.log(this.componentswe)
  },
  methods: {
    runAnsibleScript: function (script) {
      this.ajaxRequest = true
      this.loading = true
      this.script_result = ''
      this.$http.put('/ansiblescripts', script).then(function (response) {
        for (var key in response.data) {
          if (key === 'ChanID' && response.data) {
            this.doRequest(response.data.ChanID)
          } else {
            this.script_result = 'Internal Backend Error'
          }
          this.loading = false
        }
      }, function (error) {
        this.script_result = '<div class="alert alert-danger col-md-6">' + error.data.Error + '</div>'
        console.log(error)
        this.loading = false
      })
    },
    doRequest: function (chanID) {
      this.$http.get('/ansiblescript_out/' + chanID).then(function (response) {
        this.script_result = this.script_result + (response.data.cmd).replace(/\n/g, '<br/>')
        console.log(response)
        if (response.data.cmd !== 'done') {
          // setTimeout(this.doRequest(chanID),200)
          this.doRequest(chanID)
        }
      }, function (error) {
        console.log(error)
      })
    }
  }
}
</script>

<style scoped>
  body {
    background: #2B2B2B !important;
  }
  .shell {
    background: #212121;
    padding-top: 20px;
    overflow-y: auto;
    -webkit-box-shadow: inset 1px 7px 18px -5px rgba(0,0,0,0.75);
    -moz-box-shadow: inset 1px 7px 18px -5px rgba(0,0,0,0.75);
    box-shadow: inset 1px 7px 18px -5px rgba(0,0,0,0.75);
  }

  .scroll {
    overflow-y: auto;
    height:100%;
    padding-right: 0px !important;
  }

  @media screen and (min-width: 1024px) {
    shell{
      height:350px;
    }
  }
</style>
