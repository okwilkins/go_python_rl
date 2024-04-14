FROM python:3.12-alpine as base

ENV PATH="/usr/local/go/bin:$PATH"
ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONUNBUFFERED=1
# Disable CUDA support for PyTorch
ENV USE_CUDA=0

# Create user and group
RUN addgroup -S rluser && adduser -S rluser -G rluser

# Copy files in
COPY --from=golang:1.22-alpine /usr/local/go/ /usr/local/go/


FROM base as go-build

# Copy Go files
COPY go.mod /tmp/go-build/go.mod
COPY main.go /tmp/go-build/main.go
COPY go_xo /tmp/go-build/go_xo

WORKDIR /tmp/go-build

# Build Go binary
RUN go build -ldflags="-s -w" -o /usr/local/bin/go_xo


FROM base as go-base

# Copy Go binary
COPY --from=go-build /usr/local/bin/go_xo /usr/local/bin/go_xo


FROM go-base as production

# Make Python venv
RUN python -m venv /usr/local/python-venv
RUN source /usr/local/python-venv/bin/activate

USER rluser

# Install Python dependencies
COPY requirements.txt /tmp/requirements.txt
RUN pip install -r /tmp/requirements.txt

# Clean up
USER root
RUN rm -rf /tmp/*

USER rluser


FROM production as development

RUN source /usr/local/python-venv/bin/activate

USER rluser

COPY requirements_dev.txt /tmp/requirements_dev.txt
RUN pip install -r /tmp/requirements_dev.txt

# Clean up
USER root
RUN rm -rf /tmp/*

USER rluser
