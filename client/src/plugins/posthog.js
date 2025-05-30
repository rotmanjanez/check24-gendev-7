import posthog from "posthog-js";

export default {
  install(app) {
    app.config.globalProperties.$posthog = posthog.init(
      'phc_kAeue7TITO3HKu2kdNnRDo4Rby3Vc4Mb3YfOTri2re9',
      {
        api_host: 'https://eu.i.posthog.com',
      }
    );
  },
};