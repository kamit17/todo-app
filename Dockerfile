FROM jenkins/jenkins:lts

USER root

# Replace 975 with your host's Docker group GID
ARG DOCKER_GID=975


RUN apt-get update && apt-get install -y docker.io && \
    groupmod -g ${DOCKER_GID} docker && \
    usermod -aG docker jenkins

USER jenkins
