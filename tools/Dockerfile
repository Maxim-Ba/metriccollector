FROM golang:1.24.1

ENV TASK_VERSION=3.40.1
ENV UNIT_TEST_VERSION=0.11.5
ENV STATIC_TEST_VERSION=0.11.5

WORKDIR /app
COPY ./Taskfile.yml .
RUN  curl -sfL "https://github.com/Yandex-Practicum/go-autotests/releases/download/v${STATIC_TEST_VERSION}/statictest" -o /tmp/statictest && \
curl -sfL "https://github.com/Yandex-Practicum/go-autotests/releases/download/v${UNIT_TEST_VERSION}/metricstest" -o /tmp/metricstestbeta && \
curl -sfL "https://github.com/go-task/task/releases/download/v${TASK_VERSION}/task_linux_amd64.tar.gz" -o /tmp/task.tar.gz && \
tar -xzf /tmp/task.tar.gz -C /tmp && \
mv /tmp/task /usr/local/bin/ && \
mv /tmp/metricstestbeta /usr/local/bin/ && \
mv /tmp/statictest /usr/local/bin/ && \
chmod +x /usr/local/bin/* && \
    rm -rf /tmp/*

