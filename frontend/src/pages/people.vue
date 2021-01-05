<template>
  <div>
    <div v-infinite-scroll="loadMore" class="p-page p-page-albums" :infinite-scroll-disabled="scrollDisabled"
         :infinite-scroll-distance="10" :infinite-scroll-listen-for-event="'scrollRefresh'">

      <v-container v-if="loading" fluid class="pa-4">
        <v-progress-linear color="secondary-dark" :indeterminate="true"></v-progress-linear>
      </v-container>

      <v-toolbar flat color="secondary" :dense="$vuetify.breakpoint.smAndDown">
        <v-toolbar-title>
          <translate>WIP</translate>
        </v-toolbar-title>

        <v-spacer></v-spacer>
      </v-toolbar>

      <v-container>
        <p>
          Issues labeled <a href="https://github.com/photoprism/photoprism/labels/help%20wanted">help wanted</a> /
          <a href="https://github.com/photoprism/photoprism/labels/easy">easy</a> can be good (first)
          contributions.
          Our <a href="https://github.com/photoprism/photoprism/wiki">Developer Guide</a> contains all information
          necessary to get you started.
        </p>
      </v-container>
      <v-container>
        <v-layout row wrap class="p-people-results">
          <v-flex
              v-for="(people, index) in results"
              :key="index"
              :data-uid="people.UID"
              class="p-people"
              xs6 sm4 md3 lg2 d-flex
          >
            <v-hover>
              <v-card slot-scope="{ hover }" tile
                      class="accent lighten-3"
                      :to="{name: view, params: {uid: people.UID, fullname: people.FullName}}"
                      @contextmenu="onContextMenu($event, index)"
              >

                <h1
                    class="body-2 ma-0 action-title-edit"
                    :data-uid="people.UID"
                >
                  {{ people.FullName }}
                </h1>
                <p v-if="people.BoD!=null">Born at {{ people.BoD }}.</p>
                <p v-if="people.DeadDate!=null">Dead at {{ people.DeadDate }}.</p>

              </v-card>
            </v-hover>
          </v-flex>
        </v-layout>
      </v-container>
    </div>
  </div>
</template>

<script>

import People from "model/people";
import {DateTime} from "luxon";
import Event from "pubsub-js";
import RestModel from "model/rest";
import {MaxItems} from "common/clipboard";
import Notify from "common/notify";

export default {
  name: 'people',
  data() {
    const query = this.$route.query;
    const routeName = this.$route.name;
    const q = query["q"] ? query["q"] : "";
    const settings = {};

    return {
      featureShare: this.$config.feature('share'),
      subscriptions: [],
      listen: false,
      dirty: false,
      results: [],
      loading: true,
      scrollDisabled: true,
      pageSize: 24,
      offset: 0,
      page: 0,
      settings: settings,
      routeName: routeName,
      mouseDown: {
        index: -1,
        timeStamp: -1,
      },
      lastId: "",
      model: new People(),
    };
  },
  created() {
    this.search();

    //this.subscriptions.push(Event.subscribe("peoples", (ev, data) => this.onUpdate(ev, data)));

    this.subscriptions.push(Event.subscribe("touchmove.top", () => this.refresh()));
    this.subscriptions.push(Event.subscribe("touchmove.bottom", () => this.loadMore()));
  },
  destroyed() {
    for (let i = 0; i < this.subscriptions.length; i++) {
      Event.unsubscribe(this.subscriptions[i]);
    }
  },
  methods: {
    loadMore() {
      if (this.scrollDisabled) return;

      this.scrollDisabled = true;
      this.listen = false;

      const count = this.dirty ? (this.page + 2) * this.pageSize : this.pageSize;
      const offset = this.dirty ? 0 : this.offset;

      const params = {
        count: count,
        offset: offset,
      };

      Object.assign(params, this.lastFilter);

      if (this.staticFilter) {
        Object.assign(params, this.staticFilter);
      }

      People.search(params).then(response => {
        this.results = this.dirty ? response.models : this.results.concat(response.models);

        this.scrollDisabled = (response.models.length < count);

        if (this.scrollDisabled) {
          this.offset = offset;

          if (this.results.length > 1) {
            this.$notify.info(this.$gettextInterpolate(this.$gettext("All %{n} people loaded"), {n: this.results.length}));
          }
        } else {
          this.offset = offset + count;
          this.page++;

          this.$nextTick(() => {
            if (this.$root.$el.clientHeight <= window.document.documentElement.clientHeight) {
              this.$emit("scrollRefresh");
            }
          });
        }
      }).catch(() => {
        this.scrollDisabled = false;
      }).finally(() => {
        this.dirty = false;
        this.loading = false;
        this.listen = true;
      });
    },
    searchParams() {
      const params = {
        count: this.pageSize,
        offset: this.offset,
      };

      // Object.assign(params, this.filter);

      /* if (this.staticFilter) {
        Object.assign(params, this.staticFilter);
      }*/

      return params;
    },
    search() {
      this.scrollDisabled = true;

      // Don't query the same data more than once
      /*if (JSON.stringify(this.lastFilter) === JSON.stringify(this.filter)) {
        this.$nextTick(() => this.$emit("scrollRefresh"));
        return;
      }

      Object.assign(this.lastFilter, this.filter);
      */
      this.offset = 0;
      this.page = 0;
      this.loading = true;
      this.listen = false;

      const params = this.searchParams();

      People.search(params).then(response => {
        this.offset = this.pageSize;

        this.results = response.models;

        this.scrollDisabled = (response.models.length < this.pageSize);

        if (this.scrollDisabled) {
          if (!this.results.length) {
            this.$notify.warn(this.$gettext("No people found"));
          } else if (this.results.length === 1) {
            this.$notify.info(this.$gettext("One people found"));
          } else {
            this.$notify.info(this.$gettextInterpolate(this.$gettext("%{n} people found"), {n: this.results.length}));
          }
        } else {
          this.$notify.info(this.$gettext('More than 20 peoples found'));

          this.$nextTick(() => {
            if (this.$root.$el.clientHeight <= window.document.documentElement.clientHeight) {
              this.$emit("scrollRefresh");
            }
          });
        }
      }).finally(() => {
        this.dirty = false;
        this.loading = false;
        this.listen = true;
      });
    },
    refresh() {
      if (this.loading) return;
      this.loading = true;
      this.page = 0;
      this.dirty = true;
      this.scrollDisabled = false;
      this.loadMore();
    },
  }
};
</script>
