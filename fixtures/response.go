package fixtures

var DockerSwarmServicesResponse = `
[
	{
		"ID": "bs_id",
		"Version": {
			"Index": 98933
		},
		"CreatedAt": "2017-12-21T16:31:11.466402605Z",
		"UpdatedAt": "2017-12-21T16:31:53.795497001Z",
		"Spec": {
			"Name": "bs",
			"Labels": {},
			"TaskTemplate": {
				"ContainerSpec": {
          "Image": "jenkins.docker.test.com/jenkins-test:v1",
          "Env": [
              "JENKINS_ACL_MEMBERS_admin=rapid*rAPIdteam",
              "JENKINS_ACL_MEMBERS_developer=",
              "JENKINS_URL=https://bs.jenkins.test.com",
              "GHE_KEY=test_key",
              "GHE_SECRET=test_secret",
              "ADMIN_SSH_PUBKEY=ssh-rsa abcdefghijklmnop jenkins@test.com",
              "JENKINS_SLAVE_AGENT_PORT=12345",
              "JAVA_OPTS=-Xms4g -Xmx4g"
          ]
				},
				"ForceUpdate": 0,
				"Runtime": "container"
			},
			"Mode": {
				"Replicated": {
					"Replicas": 1
				}
			}
		},
		"PreviousSpec": {
			"Name": "bs",
			"Labels": {},
			"TaskTemplate": {
				"ContainerSpec": {
					"Image": "nginx:latest@sha256:test_bs_sha",
					"Env": [
						"JENKINS_URL=jenkins.test.com"
					],
					"DNSConfig": {}
				},
				"Resources": {
					"Limits": {},
					"Reservations": {}
				},
				"Placement": {
					"Platforms": []
				},
				"ForceUpdate": 0,
				"Runtime": "container"
			},
			"Mode": {
				"Replicated": {
					"Replicas": 1
				}
			},
			"EndpointSpec": {
				"Mode": "vip"
			}
		},
		"Endpoint": {
			"Spec": {}
		},
		"UpdateStatus": {
			"State": "completed",
			"StartedAt": "2017-12-21T16:31:45.393048402Z",
			"CompletedAt": "2017-12-21T16:31:53.795468535Z",
			"Message": "update completed"
		}
	},
  {
		"ID": "loatest_id",
		"Version": {
			"Index": 98933
		},
		"CreatedAt": "2017-12-21T16:31:11.466402605Z",
		"UpdatedAt": "2017-12-21T16:31:53.795497001Z",
		"Spec": {
			"Name": "loadtest",
			"Labels": {},
			"TaskTemplate": {
				"ContainerSpec": {
          "Image": "jenkins.docker.test.com/jenkins-test:v1",
          "Env": [
              "JENKINS_ACL_MEMBERS_admin=rapid*rAPIdteam",
              "JENKINS_ACL_MEMBERS_developer=",
              "JENKINS_URL=https://loadtest.jenkins.test.com",
              "GHE_KEY=test_key",
              "GHE_SECRET=test_secret",
              "ADMIN_SSH_PUBKEY=ssh-rsa abcdefghijklmnop jenkins@test.com",
              "JENKINS_SLAVE_AGENT_PORT=12345",
              "JAVA_OPTS=-Xms4g -Xmx4g"
          ]
				},
				"ForceUpdate": 0,
				"Runtime": "container"
			},
			"Mode": {
				"Replicated": {
					"Replicas": 1
				}
			}
		},
		"PreviousSpec": {
			"Name": "loadtest",
			"Labels": {},
			"TaskTemplate": {
				"ContainerSpec": {
					"Image": "nginx:latest@sha256:test_loadtest_sha",
					"Env": [
						"JENKINS_URL=jenkins.test.com"
					],
					"DNSConfig": {}
				},
				"Resources": {
					"Limits": {},
					"Reservations": {}
				},
				"Placement": {
					"Platforms": []
				},
				"ForceUpdate": 0,
				"Runtime": "container"
			},
			"Mode": {
				"Replicated": {
					"Replicas": 1
				}
			},
			"EndpointSpec": {
				"Mode": "vip"
			}
		},
		"Endpoint": {
			"Spec": {}
		},
		"UpdateStatus": {
			"State": "completed",
			"StartedAt": "2017-12-21T16:31:45.393048402Z",
			"CompletedAt": "2017-12-21T16:31:53.795468535Z",
			"Message": "update completed"
		}
	}
]
`

var DockerSwarmServiceResponse = `
{
	"ID": "loatest_id",
	"Version": {
		"Index": 98933
	},
	"CreatedAt": "2017-12-21T16:31:11.466402605Z",
	"UpdatedAt": "2017-12-21T16:31:53.795497001Z",
	"Spec": {
		"Name": "loadtest",
		"Labels": {},
		"TaskTemplate": {
			"ContainerSpec": {
				"Image": "jenkins.docker.test.com/jenkins-test:v1",
				"Env": [
						"JENKINS_ACL_MEMBERS_admin=rapid*rAPIdteam",
						"JENKINS_ACL_MEMBERS_developer=",
						"JENKINS_URL=https://loadtest.jenkins.test.com",
						"GHE_KEY=test_key",
						"GHE_SECRET=test_secret",
						"ADMIN_SSH_PUBKEY=ssh-rsa abcdefghijklmnop jenkins@test.com",
						"JENKINS_SLAVE_AGENT_PORT=12345",
						"JAVA_OPTS=-Xms4g -Xmx4g"
				]
			},
			"ForceUpdate": 0,
			"Runtime": "container"
		},
		"Mode": {
			"Replicated": {
				"Replicas": 1
			}
		}
	},
	"PreviousSpec": {
		"Name": "loadtest",
		"Labels": {},
		"TaskTemplate": {
			"ContainerSpec": {
				"Image": "nginx:latest@sha256:test_loadtest_sha",
				"Env": [
					"JENKINS_URL=jenkins.test.com"
				],
				"DNSConfig": {}
			},
			"Resources": {
				"Limits": {},
				"Reservations": {}
			},
			"Placement": {
				"Platforms": []
			},
			"ForceUpdate": 0,
			"Runtime": "container"
		},
		"Mode": {
			"Replicated": {
				"Replicas": 1
			}
		},
		"EndpointSpec": {
			"Mode": "vip"
		}
	},
	"Endpoint": {
		"Spec": {}
	},
	"UpdateStatus": {
		"State": "completed",
		"StartedAt": "2017-12-21T16:31:45.393048402Z",
		"CompletedAt": "2017-12-21T16:31:53.795468535Z",
		"Message": "update completed"
	}
}
`

var JenkinsGetAllResponse = `
["bs", "loadtest"]
`

var GitHubTeamsResponse = `
[
	{
        "name": "Admins",
        "id": 9,
        "slug": "admins",
        "description": "",
        "privacy": "secret",
        "members_count": 14,
        "repos_count": 22,
        "organization": {
            "login": "Jenkins",
            "id": 28,
            "description": "",
            "type": "Organization"
        }
    },
    {
        "name": "rAPIdteam",
        "id": 301,
        "slug": "rapidteam",
        "description": "",
        "privacy": "closed",
        "members_count": 15,
        "repos_count": 63,
        "organization": {
            "login": "rapid",
            "id": 1235,
            "description": "Cloud & Automation Services",
            "name": "RapiD",
            "type": "Organization"
        }
    }
]
`

var GitHubUserResponse = `
{
    "login": "RaphaelSantoDomingo"
}
`
