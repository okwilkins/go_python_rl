FROM python:3.12-alpine as base

ENV PATH="/usr/local/go/bin:$PATH"
ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONUNBUFFERED=1
# Disable CUDA support for PyTorch
ENV USE_CUDA=0

# Create user and group
RUN addgroup -S rluser && adduser -S rluser -G rluser

# Copy Go binaries
COPY --from=golang:1.22-alpine /usr/local/go/ /usr/local/go/


FROM base as go-build

# Copy Go files
COPY go.mod /tmp/go-build/go.mod
COPY main.go /tmp/go-build/main.go
COPY go_xo /tmp/go-build/go_xo

WORKDIR /tmp/go-build

# Build Go binary
RUN go build -ldflags="-s -w" -o /usr/local/bin/go_xo


FROM base as production

COPY --from=go-build /usr/local/bin/go_xo /usr/local/bin/go_xo

USER rluser

# Make Python venv
RUN python -m venv /home/rluser/python-venv

# Install Python dependencies
COPY requirements.txt /tmp/requirements.txt
RUN /home/rluser/python-venv/bin/python -m pip install -r /tmp/requirements.txt

# Copy Python files
COPY main.py /home/rluser/go_python_rl/main.py
COPY python_xo /home/rluser/go_python_rl/python_xo
ENV PYTHONPATH=$PYTHONPATH:/home/rluser/go_python_rl

# Clean up
USER root
RUN rm -rf /tmp/*

USER rluser


FROM production as development

COPY requirements_dev.txt /tmp/requirements_dev.txt
RUN /home/rluser/python-venv/bin/python -m pip install -r /tmp/requirements_dev.txt

# Clean up
USER root
RUN rm -rf /tmp/*

USER rluser
