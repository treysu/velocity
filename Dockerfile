FROM eclipse-temurin:21-jre-alpine

LABEL org.opencontainers.image.vendor="Dockcenter"
LABEL org.opencontainers.image.title="Velocity"
LABEL org.opencontainers.image.description="Automatically built Docker image for Velocity"
LABEL org.opencontainers.image.documentation="https://github.com/treysu/velocity/blob/main/README.md"
LABEL org.opencontainers.image.authors="Chao Tzu-Hsien <danny900714@gmail.com>"
LABEL org.opencontainers.image.licenses="MIT"

ENV JAVA_MEMORY="512M"
ENV JAVA_FLAGS="-XX:+UseStringDeduplication -XX:+UseG1GC -XX:G1HeapRegionSize=4M -XX:+UnlockExperimentalVMOptions -XX:+ParallelRefProcEnabled -XX:+AlwaysPreTouch"

WORKDIR /data

RUN mkdir -p /data && \
    apk add --upgrade --no-cache openssl && \
    addgroup -g ${VELOCITY_PGID:-1003} -S velocity && \
    adduser -u ${VELOCITY_PUID:-1002} -S velocity -G velocity && \
    chown -R velocity:velocity /data

USER velocity

VOLUME /data

EXPOSE 25577

COPY velocity/velocity-*.jar /opt/velocity/velocity.jar

ENTRYPOINT java -Xms$JAVA_MEMORY -Xmx$JAVA_MEMORY $JAVA_FLAGS -jar /opt/velocity/velocity.jar
