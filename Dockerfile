FROM rabbitmq:3.11-management

RUN apt-get update && \
apt-get install -y curl

RUN curl -LJO https://github.com/rabbitmq/rabbitmq-delayed-message-exchange/releases/download/3.11.1/rabbitmq_delayed_message_exchange-3.11.1.ez && \
mv rabbitmq_delayed_message_exchange-3.11.1.ez plugins/

RUN rabbitmq-plugins enable rabbitmq_delayed_message_exchange
RUN rabbitmq-plugins enable rabbitmq_shovel rabbitmq_shovel_management
