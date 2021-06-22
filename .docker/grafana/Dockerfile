FROM grafana/grafana:7.3.3

ARG GF_INSTALL_PLUGINS="grafana-clock-panel briangann-gauge-panel natel-plotly-panel"

ENV GF_PATHS_PROVISIONING /var/lib/grafana/provisioning

COPY /var/lib/grafana/provisioning /var/lib/grafana/provisioning

RUN echo "${GF_INSTALL_PLUGINS}" | xargs -n 1 grafana-cli plugins install
