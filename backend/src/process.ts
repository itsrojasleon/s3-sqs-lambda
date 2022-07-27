import { SQSBatchResponse, SQSEvent } from 'aws-lambda';

export const handler = async (event: SQSEvent): Promise<SQSBatchResponse> => {
  const failedMessageIds: string[] = [];

  const promises = event.Records.map(async (record) => {
    try {
      const body = JSON.parse(record.body);
      // do something with config...
      const config = body.Records[0];

      // Simulate work.
      await new Promise((resolve) => setTimeout(resolve, 1000));
    } catch (err) {
      console.log('An error ocurred: ', err);
      failedMessageIds.push(record.messageId);
    }
  });

  await Promise.allSettled(promises).then((values) => {
    console.log(values);
  });

  return {
    batchItemFailures: failedMessageIds.map((id) => ({
      itemIdentifier: id
    }))
  };
};
