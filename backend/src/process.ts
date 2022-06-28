import { APIGatewayProxyResultV2, SQSEvent } from 'aws-lambda';

export const handler = async (
  event: SQSEvent
): Promise<APIGatewayProxyResultV2> => {
  try {
    console.log('processing data...', event.Records);

    return {
      statusCode: 200,
      body: JSON.stringify({})
    };
  } catch (err: any) {
    return {
      statusCode: 500,
      body: JSON.stringify({ err })
    };
  }
};
