## Google Compute Engine
In `GCE`, you can use multiple user-data file (script or cloud-config),
through `gcloud` utility as such:
`gcloud compute instances add-metadata atomic --zone us-central1-a --metadata-from-file user-data=cloud-config.yaml`

**Note:** At the moment when providing multiple user-data files,
the order of consuming those files are non-deterministic.

## OpenStack
In 'OpenStack' you can either use, the meta-data service that
receives cloud-config files through command-line or
you can use `Config-Drive` disk images for configuration.

## Amazon Elastic Compute Cloud
You can use `Flamingo` in `EC2`, to provide meta-data or
user-data please [see](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-instance-metadata.html).
