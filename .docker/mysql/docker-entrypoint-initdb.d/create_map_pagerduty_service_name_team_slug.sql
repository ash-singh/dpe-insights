USE `dpe_insights`;

CREATE TABLE `map_pagerduty_service_name_team_slug` (
    `service_name` varchar(255) NOT NULL,
    `team_slug` varchar(255) DEFAULT NULL,
    PRIMARY KEY (`service_name`),
    UNIQUE KEY `service_name_team_slug` (`service_name`,`team_slug`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

INSERT INTO `map_pagerduty_service_name_team_slug` (`service_name`,`team_slug`) VALUES ('Account','account');
INSERT INTO `map_pagerduty_service_name_team_slug` (`service_name`,`team_slug`) VALUES ('Arbor - Critical','devops');
INSERT INTO `map_pagerduty_service_name_team_slug` (`service_name`,`team_slug`) VALUES ('Contact- Import Porcess','contact');
INSERT INTO `map_pagerduty_service_name_team_slug` (`service_name`,`team_slug`) VALUES ('Core','core');
INSERT INTO `map_pagerduty_service_name_team_slug` (`service_name`,`team_slug`) VALUES ('Datadog - Crit','devops');
INSERT INTO `map_pagerduty_service_name_team_slug` (`service_name`,`team_slug`) VALUES ('Datadog - Low','devops');
INSERT INTO `map_pagerduty_service_name_team_slug` (`service_name`,`team_slug`) VALUES ('DNS - Test','devops');
INSERT INTO `map_pagerduty_service_name_team_slug` (`service_name`,`team_slug`) VALUES ('Email Sending','Email Sending');
INSERT INTO `map_pagerduty_service_name_team_slug` (`service_name`,`team_slug`) VALUES ('Grafana APIv3','api');
INSERT INTO `map_pagerduty_service_name_team_slug` (`service_name`,`team_slug`) VALUES ('Grafana Campaign - Email','campaign');
INSERT INTO `map_pagerduty_service_name_team_slug` (`service_name`,`team_slug`) VALUES ('Marketing Automation','automation');
INSERT INTO `map_pagerduty_service_name_team_slug` (`service_name`,`team_slug`) VALUES ('MMS','devops');
INSERT INTO `map_pagerduty_service_name_team_slug` (`service_name`,`team_slug`) VALUES ('MMS - Crit','devops');
INSERT INTO `map_pagerduty_service_name_team_slug` (`service_name`,`team_slug`) VALUES ('MongoDB Cloud Manger - Production','devops');
INSERT INTO `map_pagerduty_service_name_team_slug` (`service_name`,`team_slug`) VALUES ('Redirection','sre');
INSERT INTO `map_pagerduty_service_name_team_slug` (`service_name`,`team_slug`) VALUES ('Slack oncall','devops');
INSERT INTO `map_pagerduty_service_name_team_slug` (`service_name`,`team_slug`) VALUES ('Statuscake','Statuscake');
