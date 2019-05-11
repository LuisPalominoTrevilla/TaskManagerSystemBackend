const S3FS = require('s3fs');
const officialS3 = require('aws-sdk');

const s3 = new officialS3.S3({
    region: 'us-west-1',
    accessKeyId: process.env.AWS_ACCESS_KEY,
    secretAccessKey: process.env.AWS_SECRET_KEY
})


const bucket = `${process.env.AWS_BUCKET}/tasks`;
const config = {
    region: 'us-west-1',
    accessKeyId: process.env.AWS_ACCESS_KEY,
    secretAccessKey: process.env.AWS_SECRET_KEY
}

const S3Client = new S3FS(bucket, config);

module.exports = {
    uploadFile(filename, file, extension, options) {
        return new Promise((resolve, reject) => {
            S3Client.exists(filename, exists => {
                const uniqueNumber = new Date().valueOf();
                const uniqueName = exists ? `${uniqueNumber+""}.${extension}` : filename;
                S3Client.writeFile(uniqueName, file, options)
                    .then(() => {
                        const imageUrl = `${process.env.AWS_OBJECT_PREFIX}/${bucket}/${uniqueName}`;
                        resolve(imageUrl);
                    }, (err) => {
                        reject(err);
                    });
            });
        });
    },

    deleteFile(filename) {
        const params = {
            Bucket: bucket,
            Key: filename
        };
        return new Promise((resolve, reject) => {
            s3.deleteObject(params, (err, res) => {
                if (err) reject(err);
                else resolve();
            });
        });
    }
};
