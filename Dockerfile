FROM python:3.12-alpine as base

COPY --from=golang:1.22-alpine /usr/local/go/ /usr/local/go/

# Create user and group
RUN addgroup -S rluser && adduser -S rluser -G rluser

ENV PATH="/usr/local/go/bin:$PATH"
# Disable CUDA support for PyTorch
ENV USE_CUDA=0


FROM base as production

# Make Python venv
RUN python -m venv /usr/local/python-venv
RUN source /usr/local/python-venv/bin/activate

USER rluser

# Install Python dependencies
COPY requirements.txt /tmp/requirements.txt
RUN pip install -r /tmp/requirements.txt


FROM base as development

# Make Python venv
COPY --from=production /usr/local/python-venv /usr/local/python-venv
RUN source /usr/local/python-venv/bin/activate

USER rluser

COPY requirements_dev.txt /tmp/requirements_dev.txt
RUN pip install -r /tmp/requirements_dev.txt
