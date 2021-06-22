USE `dpe_insights`;

CREATE USER 'grafanaReader' IDENTIFIED BY 'password';
 GRANT SELECT ON * TO 'grafanaReader';
